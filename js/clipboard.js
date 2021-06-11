const clipboardTest = async () => {
    const initText = await navigator.clipboard.readText();
    console.log("Initial Clipboard", initText);
    const res = await navigator.clipboard.writeText('hello, world!');
    const text = await navigator.clipboard.readText();
    console.log("Clipboard", text);
};

clipboardTest();
