import React, { useState, useEffect, useCallback } from "react";
import "./App.css";
import { FaCheck, FaExclamationTriangle } from "react-icons/fa";
import useFetch from "use-http";
import useWebSocket from "react-use-websocket";
import ReactPaginate from "react-paginate";
import Spinner from "./components/Spinner";
import { timestampToDate } from "./format/displayDate";
const hri = require("human-readable-ids").hri;
const itemsPerPage = 10;

const baseUrl = `http://${process.env.API_HOST}`;
const baseWs = `ws://${process.env.API_HOST}`;
const apiUrl = (channel) => `${baseUrl}/api/${channel || ""}`;
const wsUrl = (channel) => `${baseWs}/listen/${channel}`;

function App() {
  const [channel, setChannel] = useState(hri.random());
  const [messageHistory, setMessageHistory] = useState([]);
  const [pageCount, setPageCount] = useState(0);
  const [itemOffset, setItemOffset] = useState(0);
  const [itemCount, setItemCount] = useState(0);
  const [expandedHeaders, setExpandedHeaders] = useState(false);
  const {
    get,
    delete: clearMessages,
    response,
    loading,
  } = useFetch(baseUrl, {
    cachePolicy: "no-cache",
  });

  const loadInitialMessages = useCallback(async () => {
    const initialMessages = await get(
      `/api/${channel}?limit=${itemsPerPage}&offset=${itemOffset}`
    );
    if (response.ok) {
      setMessageHistory(
        initialMessages.data &&
          initialMessages.data
            .sort((a, b) => b.created_at - a.created_at)
            .map(({ payload,headers, failed, created_at }) => ({
              payload,
              headers,
              failed,
              created_at,
            }))
      );
      setItemCount(initialMessages.count);
    }
    //await navigator.clipboard.writeText(apiUrl(channel));
  }, [get, response, channel, setMessageHistory, itemOffset]);

  useCallback(() => {
    setPageCount(Math.ceil(itemCount / itemsPerPage));
  }, [itemCount, setPageCount]);

  useEffect(() => {
    loadInitialMessages().then(() => setSocketUrl(wsUrl(channel)))
  }, [loadInitialMessages, channel]);

  const [socketUrl, setSocketUrl] = useState(wsUrl(channel));

  const onChannelCreate = useCallback(
    (e) => {
      if (!channel) {
        return;
      }
      e.preventDefault();
      setChannel(hri.random());
      setSocketUrl(wsUrl(channel))
    },
    [setChannel, setSocketUrl, channel]
  );

  const onClear = useCallback(async () => {
    if (!itemCount) {
      return;
    }
    await clearMessages(`/api/${channel}`);
    if (response.ok) {
      setMessageHistory([]);
      setItemOffset(0);
      setPageCount(0);
      setItemCount(0);
    }
  }, [
    setMessageHistory,
    setItemCount,
    setPageCount,
    clearMessages,
    channel,
    response,
    itemCount,
  ]);

  

  useWebSocket(socketUrl, {
    onMessage: (event) => {
      // update messages for first page only
      itemOffset === 0 &&
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
              headers: message.headers,
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
    setPageCount(Math.ceil(itemCount / itemsPerPage));
  }, [itemOffset, itemCount]);

  const handlePageClick = (event) => {
    const newOffset = (event.selected * itemsPerPage) % itemCount;
    setItemOffset(newOffset);
  };

  return (
    <div className="App">
      <div className="container max-w-full py-10 flex flex-row flex-wrap justify-center gap-x-3">
        <p className="text-center w-full">Use this url or paste existing: </p>
        <div className="h-14 w-80 focus-within:text-black-600 focus:outline-none rounded-lg z-0 border-solid border-2">
          <input
            type="text"
            name="channel"
            id="channel"
            placeholder={channel}
            right={loading ? <Spinner /> : null}
            className="pl-10 pr-20 h-full w-full"
            value={channel}
            onChange={(e) => debounce(setChannel(e.target.value), 1000)}
          />
        </div>
        <span
          className="text-white mr-2 p-3 rounded-lg border-solid border-2 align-middle bg-green-500 hover:bg-green-600"
          onClick={onChannelCreate}
        >
          generate new id
        </span>
        <span
          className={`text-white mr-2 p-3 rounded-lg border-solid border-2 align-middle ${
            itemCount ? "bg-red-500 hover:bg-red-600" : "bg-gray-500"
          }`}
          onClick={onClear}
        >
          clear
        </span>
      </div>
      <div className="py-20 shadow overflow-hidden">
        {messageHistory && messageHistory.length > 0 ? (
          <table className="table-fixed border-collapse border border-blue-400 w-full">
            <thead>
              <tr className="bg-blue-600">
                <th className="border border-blue-400 w-8">
                  <span className="text-gray-100 font-semibold">Status</span>
                </th>
                <th className="border border-blue-400 w-3/6">
                  <span className="text-gray-100 font-semibold">Payload</span>
                </th>
                <th className="border border-blue-400 w-2/6">
                  <span className="text-gray-100 font-semibold">Headers</span>
                </th>
                <th className="border border-blue-400 w-12">
                  <span className="text-gray-100 font-semibold">
                    ReceivedAt
                  </span>
                </th>
              </tr>
            </thead>
            <tbody className="bg-gray-200 min-w-full">
              {messageHistory.map((message, index) => (
                <tr key={index} className="bg-white">
                  <td className="border border-blue-400">
                    <div className="self-center">
                      {message.failed === true ? (
                        <FaExclamationTriangle color="red" />
                      ) : (
                        <FaCheck color="green" />
                      )}
                    </div>
                  </td>
                  <td className="border border-blue-400 max-w-lg h-32 break-all">
                    <pre className="break-all block">
                      {prettyPrint(message.payload)}
                    </pre>
                  </td>
                  <td className="border border-blue-400 max-w-lg h-32 break-all">
                    <p onClick={() => setExpandedHeaders(!expandedHeaders)} className="text-green-600 underline">toggle headers</p>
                    <pre className={`break-all block ${!expandedHeaders && 'hidden'}`}>
                      {prettyPrint(message.headers)}
                    </pre>
                  </td>
                  <td className="border border-blue-400">
                    {timestampToDate(message.created_at)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p className="w-full text-center">
            No data
            <br />
            Send POST requests to{" "}
            <a href={apiUrl(channel)}>{apiUrl(channel)}</a>
            <br />
            Messages expire in 3 days.
          </p>
        )}
        <div id="container" className="flex flex-row justify-center">
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

const prettyPrint = (message) => {
  try {
    return JSON.stringify(JSON.parse(message), null, "\t");
  } catch (_) {
    return message;
  }
};

const debounce = (fn, wait) => {
  let timeout;
  return function() {
    const context = this;
    const args = arguments;
    const later = function() {
      timeout = null;
      fn.apply(context, args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
};

export default App;
