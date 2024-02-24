import pocketbase from "pocketbase";

export const usePb = () => {
  return new pocketbase("http://127.0.0.1:8090");
};
