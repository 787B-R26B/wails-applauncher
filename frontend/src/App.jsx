import { useState, useEffect } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { GetScriptManifest, ExecuteScript } from "../wailsjs/go/main/App";

function App() {
  const [scripts, setScripts] = useState([]);
  const [resultText, setResultText] = useState("select a scritpt to run.");

  useEffect(() => {
    GetScriptManifest()
      .then((manifestJson) => {
        try {
          const scripts = JSON.parse(manifestJson);
          setScripts(scripts);
        } catch (e) {
          console.error("failed to parse manifest:", e);
          setResultText("Error: could not load script manifest");
        }
      })
      .catch((err) => {
        console.log(err);
        setResultText("Error: could not fetch manifest");
      });
  }, []);

  function handleExecuteScript(language, filename) {
    setResultText(`Executing '${filename}'...`);
    ExecuteScript(language, filename)
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
      <h1>Wails Script Launcher</h1>
      <div className="script-list">
        {scripts.map((script, index) => (
          <div className="script-item" key={index}>
            <div className="script-info">
              <div className="script-name">{script.name}</div>
              <div className="script-desc">{script.description}</div>
            </div>
            <button
              className="btn"
              onClick={() =>
                handleExecuteScript(script.language, script.filename)
              }
            >
              Run
            </button>
          </div>
        ))}
      </div>
      <pre id="result" className="result">
        {resultText}
      </pre>
    </div>
  );
}
export default App;
