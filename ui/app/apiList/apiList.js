'use strict';
angular.module('mrredisApp')
	.constant('api', {
		/*'endPoint' :{
			'url' : 'http://10.145.208.107:8089/v1'
		},*/
		'dbStatus' : {
			'url' : '/STATUS',
			'method' : 'GET'
		},
		'dbCreate' :{
			'url' : '/CREATE',
			'method' : 'POST'
		},
		'dbDelete' :{
			'url' : '/DELETE',
			'method' : 'DELETE'
		}
	});