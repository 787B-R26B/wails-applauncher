import { Paper, Text, Code } from "@mantine/core";

interface OutputDisplayProps {
  resultText: string;
}

export function OutputDisplay({ resultText }: OutputDisplayProps) {
  return (
    <Paper withBorder p="md" radius="md" mt="lg">
      <Text size="sm" c="dimmed">
        Output
      </Text>
      <Code block>{resultText}</Code>
    </Paper>
  );
}
