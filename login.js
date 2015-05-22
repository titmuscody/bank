	function login(){
	$.ajax({
		url:"/login/",
		headers:{"Authorization":makeAuth()},
		contentType:"application/json",
		dataType:"text",
		type:"POST",
		success:function(res){
		alert("you are signed in:" + res);	
		},
		error:function(res){
			alert("failure signing in: " + res);
		}


	});
}

	function makeAuth(user, password){
		return "user:testing";
}
$(function(){
	$("#submit-button").on("click", login);
});
