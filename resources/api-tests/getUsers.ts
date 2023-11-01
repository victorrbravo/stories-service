
async function main(){
    const one_request = await fetch("http://localhost:8080/auth/login", {
        method: "POST",
        body: JSON.stringify({ email: "victorrbravo@gmail.com", password: "pass" }),
        headers: {
            "Content-Type": "application/json"
        },
    })
    const json_response = (await one_request.text()).slice(1, -1)
    console.log(json_response)
    const data = await fetch("http://localhost:8080/users", {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": json_response,
        }
    })
    const json_two = await data.text()
    console.log(JSON.parse(json_two))
}
main()
