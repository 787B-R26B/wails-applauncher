import { TextInput, Button, Group, Paper } from "@mantine/core";

interface ServerAddressFormProps {
  newServerAddress: string;
  setNewServerAddress: (address: string) => void;
  handleSaveServerAddress: () => void;
}

export function ServerAddressForm({
  newServerAddress,
  setNewServerAddress,
  handleSaveServerAddress,
}: ServerAddressFormProps) {
  return (
    <Paper withBorder p="md" radius="md">
      <Group>
        <TextInput
          label="Server Address"
          value={newServerAddress}
          onChange={(event) => setNewServerAddress(event.currentTarget.value)}
          style={{ flex: 1 }}
        />
        <Button onClick={handleSaveServerAddress} mt="lg">
          Save
        </Button>
      </Group>
    </Paper>
  );
}
