
/*使用该js需要引入jquery， Validform，layer手机提示*/
/*手机版弹出提示框*/
function alertMsg(msg){
	layer.msg(msg);
}

/*
*弹出提示信息后跳转
*/
function alertMsgUrl(msg,url,time){
	alertMsg(msg,time);
	setTimeout(function(){
		if(url){
			location.href=url;
		}
	},2000);
}

/*
*提交数据 
* obj Validform表单对象，post_url
*/

function submitForm(obj,post_url,callback){
	if(post_url == undefined || post_url == ""){
		post_url = location.href;
	}
	obj.Validform({
		tiptype:function(msg,o,cssctl){
			if(o.type != 2){
				alertMsg(msg);
			}
		
		},
		beforeSubmit:function(curform){
			jQuery.post(post_url,obj.serialize(),function(msg){
				if(typeof callback == "function"){
					callback(msg);
					return true;
				}
				if(msg.code != 2){
					alertMsg(msg.msg);
				}
				if(msg.code){
					setTimeout(function(){
						if(msg.url){
							location.href=msg.url;
						}
					},1500);
				}
			},'json')
			return false;
		},
	
	});
}



/*post提交函数*/
function postData(postUrl,jsonData){
	if(jsonData==undefined){
		$.post(postUrl,function(res){
			alertMsg(res.msg);
			 setTimeout(function(){
				if(res.url){
					location.href=res.url;
				}
			},1500);
		},"json");
	}else{
		$.post(postUrl,jsonData,function(res){
			alertMsg(res.msg);
			 setTimeout(function(){
				if(res.url){
					location.href=res.url;
				}
			},1500);
		},"json");
	}
}