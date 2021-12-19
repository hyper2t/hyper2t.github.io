(function () {
	function setUtterancesTheme() {
		let theme = window.localStorage.getItem("colorMode") || "lit";
		let msg = {
			type: "set-theme",
			theme: "icy-dark"
		};
		if (theme == "lit") {
			msg.theme = "github-light"
		}
		document.querySelector('.utterances-frame').contentWindow.postMessage(msg, 'https://utteranc.es')
	}

	document.querySelector('.color_choice').addEventListener("click", () => {
		setTimeout(setUtterancesTheme, 500);
	})
	setTimeout(setUtterancesTheme, 1000);
})();
