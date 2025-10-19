import { useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { ExecuteCommand } from "../wailsjs/go/main/App";

function App() {
  const [resultText, setResultText] = useState("click the button.");

  function handleExecute() {
    ExecuteCommand("ls -l")
      .then((result) => {
        setResultText(result);
      })
      .catch((err) => {
        setResultText(err);
      });
  }

  return (
    <div id="App">
      <img src={logo} id="logo" alt="logo" />
      <h1>wails-react-launcher</h1>
      <div className="input-box">
        <button className="btn" onClick={handleExecute}>
          Execute 'ls -l'
        </button>
      </div>
      <pre id="result" className="result">
        {resultText}
      </pre>
    </div>
  );
}
export default App;
