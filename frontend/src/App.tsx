import { Container, Title, Stack } from "@mantine/core";
import { useServerAddress } from "./hooks/useServerAddress";
import { useScriptRunner } from "./hooks/useScriptRunner";
import { ServerAddressForm } from "./components/ServerAddressForm";
import { ScriptList } from "./components/ScriptList";
import { OutputDisplay } from "./components/OutputDisplay";
import { useEffect } from "react";

function App() {
  const {
    serverAddress,
    newServerAddress,
    setNewServerAddress,
    handleSaveServerAddress,
  } = useServerAddress();

  const {
    scripts,
    resultText,
    fetchManifest,
    handleShowDetails,
    handleExecuteScript,
  } = useScriptRunner(serverAddress);

  useEffect(() => {
    fetchManifest();
  }, [fetchManifest]);

  const handleSaveAddress = () => {
    handleSaveServerAddress();
    fetchManifest();
  };

  return (
    <Container p="md">
      <Stack>
        <Title order={1} style={{ textAlign: "center" }}>
          Wails Script Launcher
        </Title>

        <ServerAddressForm
          newServerAddress={newServerAddress}
          setNewServerAddress={setNewServerAddress}
          handleSaveServerAddress={handleSaveAddress}
        />

        <ScriptList
          scripts={scripts}
          handleShowDetails={handleShowDetails}
          handleExecuteScript={handleExecuteScript}
        />

        <OutputDisplay resultText={resultText} />
      </Stack>
    </Container>
  );
}

export default App;
