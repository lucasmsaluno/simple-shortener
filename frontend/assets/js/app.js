const copyBtn = document.querySelector(".copy-btn");
const resultLink = document.querySelector(".result");
const inputLink = document.querySelector(".input-link");
const shortenForm = document.querySelector("form");
const shortBtn = document.querySelector(".short-btn");
let loader = document.querySelector(".loader");

if (!loader) {
    loader = document.createElement("span");
    loader.className = "loader";
    shortBtn.appendChild(loader);
}

inputLink.value = "";
const BASE_URL = "https://simple-shortener.up.railway.app";

const copyToClipboard = async () => {
    try {
        await navigator.clipboard.writeText(resultLink.textContent);
        alert("URL copied to clipboard!");
    } catch (err) {
        console.error("Error copying: ", err);
        alert("Failed to copy the URL.");
    }
};

const isValidURL = (url) => {
    try {
        new URL(url);
        return true;
    } catch (_) {
        return false;
    }
};

const shortenURL = async (e) => {
    e.preventDefault();

    const originalURL = inputLink.value.trim();
    if (!isValidURL(originalURL)) {
        alert("Please insert a valid URL.");
        return;
    }

    loader.style.display = "inline-block";
    shortBtn.disabled = true;

    try {
        const response = await fetch(`${BASE_URL}/api/shorten`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ original_url: originalURL }),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`API error: ${response.status} - ${errorText}`);
        }

        const data = await response.json();
        console.log(data);
        resultLink.href = data.short_url;
        resultLink.textContent = data.short_url;
        copyBtn.style.display = "flex";
    } catch (err) {
        console.error(err);
        alert("Error shortening the URL: " + err.message);
    } finally {
        loader.style.display = "none";
        shortBtn.disabled = false;
    }
};

copyBtn.addEventListener("click", copyToClipboard);
shortenForm.addEventListener("submit", shortenURL);
