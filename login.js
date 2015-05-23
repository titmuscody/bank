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
            alert(res);
            },
            error:function(res){
            //alert(res);
            }});
		alert("you are signed in:" + res);	
		},
		error:function(res){
			//alert("failure signing in: " + res);
		}


	});
}

	function makeAuth(key){
		var val = $("#username").val() + ':' + CryptoJS.HmacSHA512($("#password").val(), key);
        alert(val);
        return val.toString(CryptoJS.enc.Hex);
}
$(function(){
	$("#submit-button").on("click", login);
});
