	function login(){
	$.ajax({
		url:"/login/",
		headers:{"Authorization":$("#username").val()},
		contentType:"application/json",
		dataType:"text",
		type:"POST",
		success:function(res){
            $.ajax({
                url:"/login/",
            headers:{"Authorization":makeAuth(res)},
            contentType:"application/json",
            dataType:"text",
            type:"POST",
            success:function(res){
                if(res != "")
                    alert(res);
            },
            error:function(res){
                if(res != "")
                    alert(res);
            }});
		},
		error:function(res){
		}


	});
}

	function makeAuth(key){
		var val = $("#username").val() + ':' + CryptoJS.HmacSHA512($("#password").val(), key);
        return val.toString(CryptoJS.enc.Hex);
}
$(function(){
	$("#submit-button").on("click", login);
});
