import { Application } from "../hooks/useApplications";

export async function getApplications(
  serverAddress: string,
): Promise<Application[]> {
  const cleanedAddress = serverAddress.replace(/\/$/, "");
  const response = await fetch(`${cleanedAddress}/api/v1/applications`);
  if (!response.ok) {
    throw new Error(`Server responded with ${response.status}`);
  }
  return response.json();
}

export async function getArtifact(
  serverAddress: string,
  appName: string,
): Promise<Blob> {
  const cleanedAddress = serverAddress.replace(/\/$/, "");
  const artifactUrl = `${cleanedAddress}/api/v1/applications/${encodeURIComponent(
    appName,
  )}/artifact`;
  const response = await fetch(artifactUrl);
  if (!response.ok) {
    throw new Error(`Failed to fetch artifact: ${await response.text()}`);
  }
  return response.blob();
}
