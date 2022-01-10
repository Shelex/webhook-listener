import React, { useState, useEffect, useCallback } from "react";
import "./App.css";
import { FaRedo, FaCheck, FaExclamationTriangle } from "react-icons/fa";
import useFetch from "use-http";
import useWebSocket from "react-use-websocket";
import ReactPaginate from "react-paginate";
import Spinner from "./components/Spinner";
import { timestampToDate } from "./format/displayDate";
const hri = require("human-readable-ids").hri;
const itemsPerPage = 10;

const wsUrl = (channel) => `ws://localhost:8080/listen/${channel}`;
const apiUrl = (channel) => `http://localhost:8080/api/${channel}`;

function App() {
  const [channel, setChannel] = useState(hri.random());
  const [messageHistory, setMessageHistory] = useState([]);
  const [pageCount, setPageCount] = useState(0);
  const [itemOffset, setItemOffset] = useState(0);
  const [itemCount, setItemCount] = useState(0);
  const { get, response, loading } = useFetch(`http://localhost:8080`, {
    cachePolicy: "no-cache",
  });

  const loadInitialMessages = useCallback(async () => {
    const initialMessages = await get(
      `api/${channel}?limit=${itemsPerPage}&offset=${itemOffset}`
    );
    if (response.ok) {
      setMessageHistory(
        initialMessages.data &&
          initialMessages.data
            .sort((a, b) => b.created_at - a.created_at)
            .map(({ payload, failed, created_at }) => ({
              payload,
              failed,
              created_at,
            }))
      );
      setItemCount(initialMessages.count);
    }
    await navigator.clipboard.writeText(apiUrl(channel));
  }, [get, response, channel, setMessageHistory, itemOffset]);

  useCallback(() => {
    setPageCount(Math.ceil(itemCount / itemsPerPage));
  }, [itemCount, setPageCount]);

  const [socketUrl, setSocketUrl] = useState(wsUrl(channel));

  const onChannelCreate = useCallback(
    async (e) => {
      if (!channel) {
        return;
      }
      e.preventDefault();
      setChannel(hri.random());
    },
    [setChannel, channel]
  );

  useWebSocket(socketUrl, {
    onMessage: (event) => {
      setMessageHistory((prev) => {
        if (prev.length >= 9) {
          prev.length = 9;
        }
        let message = event.data;
        try {
          message = JSON.parse(message);
        } catch (e) {}
        return [
          {
            payload: message.payload,
            failed: !message.ok,
            created_at: Math.floor(Date.now() / 1000),
          },
          ...prev,
        ];
      });
      setItemCount(itemCount + 1);
    },
  });

  useEffect(() => {
    loadInitialMessages().then(() => setSocketUrl(wsUrl(channel)));
  }, [loadInitialMessages, setSocketUrl, channel]);

  useEffect(() => {
    // Fetch items from another resources.
    const endOffset = itemOffset + itemsPerPage;
    console.log(`Loading items from ${itemOffset} to ${endOffset}`);
    setPageCount(Math.ceil(itemCount / itemsPerPage));
  }, [itemOffset, itemCount]);

  // Invoke when user click to request another page.
  const handlePageClick = (event) => {
    const newOffset = (event.selected * itemsPerPage) % itemCount;
    console.log(
      `User requested page number ${event.selected}, which is offset ${newOffset}`
    );
    setItemOffset(newOffset);
  };

  return (
    <div className="App">
      <header className="App-header"></header>
      <div className="container grid justify-items-center">
        <p>Use this url or paste existing</p>
        <label htmlFor="channel">
          <div className="h-14 w-96 text-gray-400 focus-within:text-black-600 focus:outline-none rounded-lg z-0 border-solid border-2 flex items-center relative">
            <input
              type="text"
              name="channel"
              id="channel"
              placeholder={channel}
              right={loading ? <Spinner /> : null}
              className="pl-10 pr-20 h-full w-full"
              value={channel}
              onChange={(e) => e.target.value && setChannel(e.target.value)}
            />
            <div className="absolute right-0" onClick={onChannelCreate}>
              <FaRedo
                color="white"
                className="h-10 w-20 right-0 text-white rounded-lg bg-green-500 hover:bg-green-600 p-2"
              />
            </div>
          </div>
        </label>
      </div>
      <div className="max-w-6xl py-20 px-4 mx-auto grid justify-items-center">
        {messageHistory && messageHistory.length > 0 ? (
          <table className="table-auto border-collapse border border-blue-400">
            <thead className="space-x-1">
              <tr className="bg-blue-600 px-auto py-auto">
                <th className="max-w-xs border border-blue-400">
                  <span className="text-gray-100 font-semibold">Status</span>
                </th>
                <th className="w-1/2 min-w-[50%] text-base max-w-prose border border-blue-400">
                  <span className="text-gray-100 font-semibold">Payload</span>
                </th>
                <th className="max-w-max border border-blue-400">
                  <span className="text-gray-100 font-semibold">
                    ReceivedAt
                  </span>
                </th>
              </tr>
            </thead>
            <tbody className="bg-gray-200">
              {messageHistory.map((message, index) => (
                <tr key={index} className="bg-white">
                  <td className="border border-blue-400">
                    {message.failed === true ? (
                      <FaExclamationTriangle
                        color="red"
                        textDecoration="copy to clibpoard"
                      />
                    ) : (
                      <FaCheck color="green" />
                    )}
                  </td>
                  <td className="border border-blue-400">
                    <pre>{message.payload}</pre>
                  </td>
                  <td className="border border-blue-400">
                    {timestampToDate(message.created_at)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p>No data</p>
        )}
        <div id="container" className="flex">
          <ReactPaginate
            nextLabel="next >"
            onPageChange={handlePageClick}
            pageRangeDisplayed={3}
            marginPagesDisplayed={2}
            pageCount={pageCount}
            previousLabel="< previous"
            pageClassName="page-item"
            pageLinkClassName="page-link"
            previousClassName="page-item"
            previousLinkClassName="page-link"
            nextClassName="page-item"
            nextLinkClassName="page-link"
            breakLabel="..."
            breakClassName="page-item"
            breakLinkClassName="page-link"
            containerClassName="pagination"
            activeClassName="page-active"
            disabledClassName="page-disabled"
            renderOnZeroPageCount={null}
          />
        </div>
      </div>
    </div>
  );
}

export default App;
