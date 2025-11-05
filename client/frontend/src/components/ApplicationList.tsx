import { Card, Text, Button, Group, Stack } from "@mantine/core";
import { Application } from "../hooks/useApplications"; // Import the shared interface

interface ApplicationListProps {
  applications: Application[];
  handleShowDetails: (description: string) => void;
  handleDownloadArtifact: (name: string) => void;
}

export function ApplicationList({
  applications,
  handleShowDetails,
  handleDownloadArtifact,
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
              <Button onClick={() => handleDownloadArtifact(app.name)}>
                Download
              </Button>
            </Group>
          </Group>
        </Card>
      ))}
    </Stack>
  );
}
