import { Card, Text, Button, Group, Stack } from "@mantine/core";
import { Application } from "../hooks/useApplications";

interface ApplicationListProps {
  applications: Application[];
  handleShowDetails: (description: string) => void;
  handleRunApplication: (app: Application) => void;
}

export function ApplicationList({
  applications,
  handleShowDetails,
  handleRunApplication,
}: ApplicationListProps) {
  return (
    <Stack gap="sm">
      {applications.map((app, index) => (
        <Card shadow="sm" padding="lg" radius="md" withBorder key={index}>
          <Group justify="space-between">
            <Text fw={500}>{app.name}</Text>
            <Group>
              <Button
                variant="light"
                onClick={() => handleShowDetails(app.description)}
              >
                Details
              </Button>
              <Button onClick={() => handleRunApplication(app)}>Run</Button>
            </Group>
          </Group>
        </Card>
      ))}
    </Stack>
  );
}
