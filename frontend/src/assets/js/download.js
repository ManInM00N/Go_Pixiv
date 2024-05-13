import axios from "axios";

function Download(){
    axios.post("http://127.0.0.1:7234/api/download",{
        type: "Pid",

    },{
        headers:{
            'Content-Type': 'application/json'
        }
    }).then(response=>{

    }).catch(error=>{
        console.error
    })
}