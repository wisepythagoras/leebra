const fetchTest = async () => {
    console.log('Start HTTP request');

    const response = await fetch('https://example.com/');

    console.log(await response.text());
};

fetchTest();
