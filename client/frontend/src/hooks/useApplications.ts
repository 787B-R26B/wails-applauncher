import { useState, useCallback } from "react";
import { notifications } from "@mantine/notifications";
import { SaveAndRunArtifact } from "../../wailsjs/go/main/App";

export interface Application {
  name: string;
  description: string;
  artifact_type: string;
  build_command: string;
  artifact_path: string;
  run_command: string;
}

export function useApplications(serverAddress: string) {
  const [applications, setApplications] = useState<Application[]>([]);
  const [resultText, setResultText] = useState<string>(
    "Select an application to run, or view details",
  );

  const fetchApplications = useCallback(() => {
    if (!serverAddress) return;
    const cleanedAddress = serverAddress.replace(/\/$/, "");
    fetch(`${cleanedAddress}/api/v1/applications`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Server responded with ${response.status}`);
        }
        return response.json();
      })
      .then((data: Application[]) => {
        setApplications(data);
      })
      .catch((err) => {
        console.error("Failed to fetch applications:", err);
        notifications.show({
          title: "Error",
          message: `Could not fetch applications: ${err.message}`,
          color: "red",
        });
      });
  }, [serverAddress]);

  const handleShowDetails = (description: string) => {
    setResultText(description);
  };

  const handleRunApplication = async (app: Application) => {
    if (!serverAddress) {
      notifications.show({
        title: "Error",
        message: "Server address is not set.",
        color: "red",
      });
      return;
    }
    setResultText(`Preparing to run '${app.name}'...`);

    try {
      // 1. Fetch artifact from server
      const cleanedAddress = serverAddress.replace(/\/$/, "");
      const artifactUrl = `${cleanedAddress}/api/v1/applications/${encodeURIComponent(
        app.name,
      )}/artifact`;
      const response = await fetch(artifactUrl);
      if (!response.ok) {
        throw new Error(`Failed to fetch artifact: ${await response.text()}`);
      }
      const blob = await response.blob();
      const fileData = new Uint8Array(await blob.arrayBuffer());

      // 2. Pass data to client's Go backend to save and run
      setResultText(`Executing '${app.name}'...`);
      const output = await SaveAndRunArtifact(
        app.artifact_type === "zip",
        Array.from(fileData), // Convert Uint8Array to a plain array for Wails
        app.run_command,
      );

      // 3. Display result
      setResultText(output);
      notifications.show({
        title: "Execution Finished",
        message: `'${app.name}' finished running.`,
        color: "green",
      });
    } catch (err: any) {
      console.error("Failed to run application:", err);
      setResultText(err.message);
      notifications.show({
        title: "Error",
        message: err.message,
        color: "red",
      });
    }
  };

  return {
    applications,
    resultText,
    fetchApplications,
    handleShowDetails,
    handleRunApplication,
  };
}
