import { useState, useEffect } from "react";
import {
  GetScriptManifest,
  ExecuteScriptInTerminal,
  GetServerAddress,
  SetServerAddress,
} from "../wailsjs/go/main/App";
import {
  Container,
  Title,
  TextInput,
  Button,
  Card,
  Text,
  Group,
  Code,
  Stack,
  useMantineTheme,
  Paper,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";

// Define the type for a single script
interface Script {
  name: string;
  description: string;
  language: string;
  filename: string;
}

function App() {
  const [scripts, setScripts] = useState<Script[]>([]);
  const [resultText, setResultText] = useState<string>(
    "Select a script to run, or view details",
  );
  const [serverAddress, setServerAddress] = useState<string>("");
  const [newServerAddress, setNewServerAddress] = useState<string>("");
  const theme = useMantineTheme();

  useEffect(() => {
    function fetchManifest() {
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
    }

    GetServerAddress().then((address) => {
      setServerAddress(address);
      setNewServerAddress(address);
      fetchManifest();
    });
  }, []);

  function handleShowDetails(description: string) {
    setResultText(description);
  }

  function handleExecuteScript(language: string, filename: string) {
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
  }

  function handleSaveServerAddress() {
    SetServerAddress(newServerAddress)
      .then(() => {
        setServerAddress(newServerAddress);
        notifications.show({
          title: "Success",
          message: "Server address updated successfully.",
          color: "green",
        });
        // Refetch manifest after address change
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
      })
      .catch((err) => {
        notifications.show({
          title: "Error",
          message: `Error updating server address: ${err}`,
          color: "red",
        });
      });
  }

  return (
    <Container p="md">
      <Stack>
        <Title order={1} style={{ textAlign: "center" }}>
          Wails Script Launcher
        </Title>

        <Paper withBorder p="md" radius="md">
          <Group>
            <TextInput
              label="Server Address"
              value={newServerAddress}
              onChange={(event) =>
                setNewServerAddress(event.currentTarget.value)
              }
              style={{ flex: 1 }}
            />
            <Button onClick={handleSaveServerAddress} mt="lg">
              Save
            </Button>
          </Group>
        </Paper>

        <Stack gap="sm">
          {scripts.map((script, index) => (
            <Card shadow="sm" padding="lg" radius="md" withBorder key={index}>
              <Group justify="space-between">
                <Text fw={500}>{script.name}</Text>
                <Group>
                  <Button
                    variant="light"
                    onClick={() => handleShowDetails(script.description)}
                  >
                    Details
                  </Button>
                  <Button
                    onClick={() =>
                      handleExecuteScript(script.language, script.filename)
                    }
                  >
                    Run
                  </Button>
                </Group>
              </Group>
            </Card>
          ))}
        </Stack>

        <Paper withBorder p="md" radius="md" mt="lg">
          <Text size="sm" c="dimmed">
            Output
          </Text>
          <Code block>{resultText}</Code>
        </Paper>
      </Stack>
    </Container>
  );
}
export default App;
