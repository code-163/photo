console.log(localStorage.getItem('pid'));

//检测是否已登录状态
var agenttoken = localStorage.getItem('agenttoken');
var agentcode = localStorage.getItem('agentcode');
var aid = localStorage.getItem('aid');
var account = localStorage.getItem('account');
//alert(usercode+'=>'+usertoken+'=>'+uid+'=>'+account)
var pagename = window.location.pathname;

if (pagename == "/page/login.html") {
    if (agenttoken != null && agentcode != null) {
        $.post(server_url + 'agent/checkAgentlogin', {agentcode: agentcode, agenttoken: agenttoken}, function (res) {
            if (res.code === 0) {
                location.replace('../../../');
            } else {
                localStorage.removeItem("agentcode")
                localStorage.removeItem("agenttoken")
                localStorage.removeItem("aid")
                localStorage.removeItem("account")
                localStorage.removeItem("coder");
                localStorage.removeItem("pid");
                //layer.msg(res.msg, {icon: 2, anim: 6});
            }
        }, 'json');
    }
} else {
    if (agenttoken != null && agentcode != null) {
        $.post(server_url + 'agent/checkAgentlogin', {agentcode: agentcode, agenttoken: agenttoken}, function (res) {
            if (res.code === 0) {
                console.log("目前登录状态");
            } else {
                localStorage.removeItem("agentcode")
                localStorage.removeItem("agenttoken")
                localStorage.removeItem("aid")
                localStorage.removeItem("account")
                localStorage.removeItem("coder");
                localStorage.removeItem("pid");
                top.location.replace('/page/login.html');
                //layer.msg(res.msg, {icon: 2, anim: 6});
            }
        }, 'json');
    } else {
        top.location.replace('/page/login.html');
    }
}



