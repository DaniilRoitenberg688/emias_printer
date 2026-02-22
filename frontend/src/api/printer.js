import axios from "axios";

let apiUrl = import.meta.env.VITE_API_URL;
if (apiUrl === undefined) {
  apiUrl = "http://localhost:8000/api/v1";
  console.error("Cannot access main api, because url was not provided");
}

export const api = axios.create({
  headers: {
    "Content-Type": "application/json",
  },
});

export async function getPrinter() {
  try {
    const response = await api.get(`${apiUrl}/printer/find`);
    return response.data;
  } catch (error) {
    console.error("Error getting printer", error);
    throw error;
  }
  return null;
}

export async function printOnPrinter(printerIp) {
  try {
    let data = {
      ip: printerIp,
      text: "Hello world",
    };
    console.log("data ihfujsdfh", data);
    const response = await api.post(`${apiUrl}/printer/print`, data);
    return response.data;
  } catch (error) {
    console.error("Error printing on printer", error);
    throw error;
  }
}

export async function checkPrinter(printerIp) {
  try {
    let data = {
      ip: printerIp,
    };
    const response = await api.post(`${apiUrl}/printer/check`, data);
    return response.data;
  } catch (error) {
    console.error("Error checking printer", error);
    throw error;
  }
}
