import { useState, useCallback } from "react";
import { notifications } from "@mantine/notifications";

export interface Application {
  name: string;
  description: string;
  artifact_type: string;
  build_command: string;
  artifact_path: string;
}

export function useApplications(serverAddress: string) {
  const [applications, setApplications] = useState<Application[]>([]);
  const [resultText, setResultText] = useState<string>(
    "Select an application to download, or view details",
  );

  const fetchApplications = useCallback(() => {
    if (!serverAddress) return;
    fetch(`${serverAddress}/api/v1/applications`)
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

  const handleDownloadArtifact = (name: string) => {
    if (!serverAddress) {
      notifications.show({
        title: "Error",
        message: "Server address is not set.",
        color: "red",
      });
      return;
    }
    setResultText(`Downloading artifact for '${name}'...`);
    const url = `${serverAddress}/api/v1/applications/${encodeURIComponent(
      name,
    )}/artifact`;

    // Create a temporary anchor element to trigger the download
    const link = document.createElement("a");
    link.href = url;

    // The 'download' attribute is not strictly necessary if Content-Disposition is set,
    // but it can help.
    // link.setAttribute('download', ''); // Optional

    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    notifications.show({
      title: "Download Started",
      message: `Your download for "${name}" should now be starting.`,
      color: "green",
    });
  };

  return {
    applications,
    resultText,
    fetchApplications,
    handleShowDetails,
    handleDownloadArtifact,
  };
}
