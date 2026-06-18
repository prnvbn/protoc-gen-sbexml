async function instantiate(url, importObject) {
  const response = await fetch(url);

  if (WebAssembly.instantiateStreaming) {
    try {
      return await WebAssembly.instantiateStreaming(response, importObject);
    } catch (error) {
      if (response.headers.get("Content-Type") === "application/wasm") {
        throw error;
      }
    }
  }

  const bytes = await response.arrayBuffer();
  return WebAssembly.instantiate(bytes, importObject);
}

export const wasmReady = (async () => {
  const go = new Go();
  const result = await instantiate("main.wasm", go.importObject);
  go.run(result.instance);
})();
