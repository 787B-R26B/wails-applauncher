import { TextInput, Button, Group, Card } from "@mantine/core";

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
    <Card shadow="sm" padding="lg" radius="md" withBorder>
      <Group>
        <TextInput
          label="Server Address"
          placeholder="http://localhost:8080"
          value={newServerAddress}
          onChange={(event) => setNewServerAddress(event.currentTarget.value)}
          style={{ flex: 1 }}
        />
        <Button
          onClick={handleSaveServerAddress}
          style={{ alignSelf: "flex-end" }}
        >
          Save
        </Button>
      </Group>
    </Card>
  );
}
