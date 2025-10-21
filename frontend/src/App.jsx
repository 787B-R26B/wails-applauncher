import { useState, useEffect } from "react";
import "./App.css";
import {
  GetScriptManifest,
  ExecuteScriptInTerminal,
  GetServerAddress,
  SetServerAddress,
} from "../wailsjs/go/main/App";

function App() {
  const [scripts, setScripts] = useState([]);
  const [resultText, setResultText] = useState(
    "select a script to run, or view details",
  );
  const [serverAddress, setServerAddress] = useState("");
  const [newServerAddress, setNewServerAddress] = useState("");

  useEffect(() => {
    function fetchManifest() {
      GetScriptManifest()
        .then((manifestJson) => {
          try {
            setScripts(JSON.parse(manifestJson));
          } catch (e) {
            console.error("failed to parse manifest:", e);
            setResultText("Error: could not load script manifest");
          }
        })
        .catch((err) => {
          console.error("Failed to fetch manifest:", err);
          setResultText(
            "Error: could not fetch manifest. Is the server running?",
          );
        });
    }

    GetServerAddress().then((address) => {
      setServerAddress(address);
      setNewServerAddress(address);
      fetchManifest();
    });
  }, []);

  function handleShowDetails(description) {
    setResultText(description);
  }

  function handleExecuteScript(language, filename) {
    setResultText(`Executing '${filename}'...`);
    ExecuteScriptInTerminal(language, filename)
      .then(() => {
        setResultText("Script executed successfully");
      })
      .catch((err) => {
        setResultText(err);
      });
  }

  function handleSaveServerAddress() {
    SetServerAddress(newServerAddress)
      .then(() => {
        setServerAddress(newServerAddress);
        setResultText("Server address updated successfully.");
        // Refetch manifest after address change
        GetScriptManifest()
          .then((manifestJson) => {
            try {
              setScripts(JSON.parse(manifestJson));
            } catch (e) {
              console.error("failed to parse manifest:", e);
              setResultText("Error: could not load script manifest");
            }
          })
          .catch((err) => {
            console.error("Failed to fetch manifest:", err);
            setResultText(
              "Error: could not fetch manifest. Is the server running?",
            );
          });
      })
      .catch((err) => {
        setResultText(`Error updating server address: ${err}`);
      });
  }

  return (
    <div id="App">
      <h1>Wails Script Launcher</h1>
      <div className="server-config">
        <label htmlFor="server-address">Server Address:</label>
        <input
          id="server-address"
          type="text"
          value={newServerAddress}
          onChange={(e) => setNewServerAddress(e.target.value)}
        />
        <button className="btn" onClick={handleSaveServerAddress}>
          Save
        </button>
      </div>
      <div className="script-list">
        {scripts.map((script, index) => (
          <div className="script-item" key={index}>
            <div className="script-info">
              <div className="script-name">{script.name}</div>
            </div>
            <div className="button-group">
              <button
                className="btn"
                onClick={() => handleShowDetails(script.description)}
              >
                Details
              </button>
              <button
                className="btn"
                onClick={() =>
                  handleExecuteScript(script.language, script.filename)
                }
              >
                Run
              </button>
            </div>
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
