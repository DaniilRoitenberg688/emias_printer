import { Button, message } from "antd";
import { checkPrinter, getPrinter, printOnPrinter } from "../api/printer";

async function setPrinterIp() {
  try {
    let data = await getPrinter();
    console.log("data", data);
    if (data.ip !== undefined) {
      localStorage.setItem("printerIp", data.ip[0]);
      return true;
    }
  } catch (error) {
    console.error("Error getting printer", error);
    return false;
  }
}

const ClickButton = () => {
  const [messageApi, contextHolder] = message.useMessage();
  const throwErrorAlert = (text) => {
    messageApi.open({
      type: "error",
      content: text,
    });
  };

  const throwSuccessAlert = (text) => {
    messageApi.open({
      type: "success",
      content: text,
    });
  };

  // const [isError, setIsError] = useContext(false);
  const handleClick = async () => {
    let printer = localStorage.getItem("printerIp");
    if (printer === null || printer === undefined) {
      let res = await setPrinterIp();
      if (!res) {
        throwErrorAlert("не могу найти принтер в вашей сети");
        localStorage.removeItem("printerIp");
        return;
      }
    } else {
      let checkPrinterRes = await checkPrinter(printer);
      if (!checkPrinterRes.available) {
        let res = await setPrinterIp();
        if (!res) {
          throwErrorAlert("не могу найти принтер в вашей сети");
          localStorage.removeItem("printerIp");
          return;
        }
      }
    }
    try {
      printer = localStorage.getItem("printerIp");
      console.log("printer", printer);
      let res = await printOnPrinter(printer);
      if (res.result === "ok") {
        throwSuccessAlert("печать успешно произведена");
      }
    } catch (error) {
      console.error("Error printing", error);
      throwErrorAlert("не могу передать запрос на печать");
    }
  };

  return (
    <>
      {contextHolder}
      <Button type="primary" onClick={handleClick}>
        Print
      </Button>
    </>
  );
};

export default ClickButton;
