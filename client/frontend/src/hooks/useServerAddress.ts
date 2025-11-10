import { useState, useEffect, useCallback } from "react";
import { notifications } from "@mantine/notifications";

const SERVER_ADDRESS_STORAGE_KEY = "serverAddress";

export function useServerAddress() {
  const [serverAddress, setServerAddress] = useState<string>("");
  const [newServerAddress, setNewServerAddress] = useState<string>("");

  useEffect(() => {
    // Load the saved address from localStorage on initial render
    const savedAddress =
      localStorage.getItem(SERVER_ADDRESS_STORAGE_KEY) ||
      "http://localhost:8080";
    setServerAddress(savedAddress);
    setNewServerAddress(savedAddress);
  }, []);

  const handleSaveServerAddress = useCallback(() => {
    try {
      localStorage.setItem(SERVER_ADDRESS_STORAGE_KEY, newServerAddress);
      setServerAddress(newServerAddress);
      notifications.show({
        title: "Success",
        message: "Server address updated successfully.",
        color: "green",
      });
    } catch (err) {
      notifications.show({
        title: "Error",
        message: `Error updating server address: ${err}`,
        color: "red",
      });
    }
  }, [newServerAddress]);

  return {
    serverAddress,
    newServerAddress,
    setNewServerAddress,
    handleSaveServerAddress,
  };
}
