import { Card, Text, Button, Group, Stack } from "@mantine/core";

interface Script {
  name: string;
  description: string;
  language: string;
  filename: string;
}

interface ScriptListProps {
  scripts: Script[];
  handleShowDetails: (description: string) => void;
  handleExecuteScript: (language: string, filename: string) => void;
}

export function ScriptList({
  scripts,
  handleShowDetails,
  handleExecuteScript,
}: ScriptListProps) {
  return (
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
  );
}
