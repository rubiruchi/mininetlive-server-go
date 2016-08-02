var config = {
    host: "http://127.0.0.1:8081"
}

var mininet = {};

function ajax(method, path, data, success, error){
    $.ajax({
        url: path,
        method: method,
        contentType: "application/x-www-form-urlencoded",
        data: data,
        success: function(rsp){
            if (rsp.ret == -1){
                var next = rsp.data.redirectPath;
                window.location.href = "/login.html";
                // window.location.href = "/login.html?next=" + next;
            } else {
                success(rsp);
            }
        },
        error: error
    })
}


function ajaxFile(method, path, data, success, fail, contentType){
    // contentType = contentType || "application/x-www-form-urlencoded";
    debugger
    $.ajax({
        url: path,
        method: method,
        contentType: false,
        processData: false,
        data: data,
        success: function(rsp){
            if (rsp.ret == -1){
                var next = rsp.data.redirectPath;
                window.location.href = "/login.html?next=" + next;
            } else {
                success(rsp);
            }
        },
        fail: fail
    })
}


function formatGender(gender){
    switch(gender){
        case 0:
            return "女"
        case 1:
            return "男"
        default:
            return "未知"
    }
}

function formatPlat(plat){
    switch(plat){
        case "QQ":
            return "QQ"
        case "Wechat":
            return "微信"
        case "SinaWeibo":
            return "新浪微博"
        default:
            return plat;
    }
}

function formatChannel(plat){
    switch(plat){
        case "wx":
            return "微信"
        default:
            return plat;
    }
}

function formatStreamType(type){
    switch(type){
        case 1:
            return "点播"
        case 0:
            return "直播"
        default:
            return type;
    }
}

function formatActivityState(state){
    switch(state){
        case 0:
            return "未开播";
        case 1:
            return "直播中";
        case 2:
            return "已结束"
    }
}

function changeTow(number){
    number += "";
    if (number.length == 1){
        return "0" + number;
    }
    return number;
}

function formatDateTime(date){
    date = new Date(date);
    return  date.getFullYear() + "-" + (date.getMonth() + 1) + "-" + date.getDate() + " " + date.getHours() + ":" + changeTow(date.getMinutes());
}

function renderHtmlNavbar(route){
    var $siderbar = $("#sidebar-nav");
    var navbar = '<ul id="dashboard-menu">' +
                    '<li class="index"><a href="/"><i class="icon-home"></i><span>首页</span></a></li>' +
                    '<li class="user"><a href="/user-list.html"><i class="icon-group"></i><span>用户</span></a></li>' +      
                    '<li class="order"><a href="/order-list.html"><i class="icon-signal"></i><span>收入</span></a></li>' +
                    '<li class="activity"><a href="/activity-list.html"><i class="icon-th-large"></i><span>活动</span></a></li>' +
                    '<li class="logout"><a href="#"><i class="icon-share-alt"></i><span>退出</span></a></li>' +
                    '</ul>';
    var pointer = '<div class="pointer"><div class="arrow"></div><div class="arrow_border"></div></div>';
    $siderbar.html(navbar);

    if (route){
        var current = $("#dashboard-menu ." + route);
        current.addClass('active');
        current.prepend(pointer);
    }

    // 退出
    $siderbar.on('click', '.logout', function(e){
        debugger
        // e.preventDefault();
        ajax("post", "/logout", {}, function(rsp){
            debugger
            if (rsp.ret == 0){
                window.location.href = "/login.html";
            }
        }, function(){
            debugger
        })
    })
}

function renderHtmlPagination(total, current, pageSize){
    total = parseInt(total);
    current = parseInt(current || 1);
    pageSize = parseInt(pageSize);
    
    var params = {
        pageSize: pageSize || 10
    };

    var html = "<ul>";

    // 前一页
    params.pageIndex = current - 1 || 1;
    html += '<li><a href="' + newLocationPath(params) + '">&#8249;</a></li>';

    
    for (var i = 1; i <= total; i++){
        params.pageIndex = i;
        if (current != i){
            html += '<li><a href="' + newLocationPath(params) + '">' + i + '</a></li>';
        } else {
            html += '<li><a href="' + newLocationPath(params) + '" class="active" >' + i + '</a></li>';
        }   
    }

    // 后一页
    params.pageIndex = current + 1;
    if (params.pageIndex > total){
        params.pageIndex = total;
    }
    html += '<li><a href="' + newLocationPath(params) + '">&#8250;</a></li>';

    html += '</ul>';
    return html
}


function newLocationPath(params){
    return location.pathname + "?" + _.stringifyUrlParams(params);
}

function formatActivityState(state){
    switch(state){
        case 0:
            return "未开播"
        case 1:
            return "直播中"
        case 2:
            return "直播结束"
        default:
            return state;
    }
}

function formateActivityType(state){
    switch(state){
        case 0:
            return "收费"
        case 1:
            return "免费"
        default:
            return state;
    }
}

function formateAppoinState(state){
    switch(state){
        case 0:
            return "未预约"
        case 1:
            return "已预约"
        default:
            return state;
    }
}

function formatePayState(state){
    switch(state){
        case 0:
            return "未支付"
        case 1:
            return "已支付"
        default:
            return state;
    }
}

mininet.ajax = ajax;
mininet.ajaxFile = ajaxFile;
mininet.formatGender = formatGender;
mininet.formatPlat = formatPlat;
mininet.formatChannel = formatChannel;
mininet.renderHtmlNavbar = renderHtmlNavbar;
mininet.renderHtmlPagination = renderHtmlPagination;
mininet.formatStreamType = formatStreamType;
mininet.formatDateTime = formatDateTime;
mininet.formatActivityState = formatActivityState;
mininet.formateActivityType = formateActivityType;
mininet.formateAppoinState = formateAppoinState;
mininet.formatePayState = formatePayState;