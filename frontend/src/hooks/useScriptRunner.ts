import { useState, useCallback } from "react";
import {
  GetScriptManifest,
  ExecuteScriptInTerminal,
} from "../../wailsjs/go/main/App";
import { notifications } from "@mantine/notifications";

interface Script {
  name: string;
  description: string;
  language: string;
  filename: string;
}

export function useScriptRunner(serverAddress: string) {
  const [scripts, setScripts] = useState<Script[]>([]);
  const [resultText, setResultText] = useState<string>(
    "Select a script to run, or view details",
  );

  const fetchManifest = () => {
    GetScriptManifest()
      .then((manifestJson) => {
        try {
          setScripts(JSON.parse(manifestJson));
        } catch (e) {
          console.error("failed to parse manifest:", e);
          notifications.show({
            title: "Error",
            message: "Could not load script manifest",
            color: "red",
          });
        }
      })
      .catch((err) => {
        console.error("Failed to fetch manifest:", err);
        notifications.show({
          title: "Error",
          message: "Could not fetch manifest. Is the server running?",
          color: "red",
        });
      });
  };

  const handleShowDetails = (description: string) => {
    setResultText(description);
  };

  const handleExecuteScript = (language: string, filename: string) => {
    setResultText(`Executing '${filename}'...`);
    ExecuteScriptInTerminal(language, filename)
      .then(() => {
        setResultText("Script executed successfully");
        notifications.show({
          title: "Success",
          message: "Script executed successfully",
          color: "green",
        });
      })
      .catch((err) => {
        setResultText(err);
        notifications.show({
          title: "Error",
          message: err,
          color: "red",
        });
      });
  };

  return {
    scripts,
    resultText,
    fetchManifest,
    handleShowDetails,
    handleExecuteScript,
  };
}
