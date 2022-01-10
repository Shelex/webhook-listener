import dayjs from "dayjs";

export const timestampToDate = (timestamp) => {
  return dayjs(timestamp * 1000).format("DD-MM-YYYY HH:mm:ss");
};
