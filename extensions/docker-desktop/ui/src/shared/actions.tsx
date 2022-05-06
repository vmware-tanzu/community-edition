export async function openBrowser(url: string) {
  return window.ddClient.host.openExternal(url);
}
