const express = require('express')
const app = express()

app.get("/", (req, res) => {
    res.send("Hello World");
});

const run = () => {
    app.listen(8080, () => console.log("Server up"))
}

run()