import { useState, useEffect, useCallback } from "react";
import { GetServerAddress, SetServerAddress } from "../../wailsjs/go/main/App";
import { notifications } from "@mantine/notifications";

export function useServerAddress() {
  const [serverAddress, setServerAddress] = useState<string>("");
  const [newServerAddress, setNewServerAddress] = useState<string>("");

  useEffect(() => {
    GetServerAddress().then((address) => {
      setServerAddress(address);
      setNewServerAddress(address);
    });
  }, []);

  const handleSaveServerAddress = useCallback(() => {
    SetServerAddress(newServerAddress)
      .then(() => {
        setServerAddress(newServerAddress);
        notifications.show({
          title: "Success",
          message: "Server address updated successfully.",
          color: "green",
        });
      })
      .catch((err) => {
        notifications.show({
          title: "Error",
          message: `Error updating server address: ${err}`,
          color: "red",
        });
      });
  }, [newServerAddress]);

  return {
    serverAddress,
    newServerAddress,
    setNewServerAddress,
    handleSaveServerAddress,
  };
}
