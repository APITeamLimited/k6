"use strict";

var HTTPResponse = ***REMOVED***
	json: function() ***REMOVED*** return JSON.parse(this.body); ***REMOVED***,
***REMOVED***

__internal__.modules.http.get = function() ***REMOVED***
	return __internal__.modules.http.do.apply(this, _.concat(['GET'], arguments));
***REMOVED***

__internal__.modules.http.post = function() ***REMOVED***
	return __internal__.modules.http.do.apply(this, _.concat(['POST'], arguments));
***REMOVED***

__internal__.modules.http.put = function() ***REMOVED***
	return __internal__.modules.http.do.apply(this, _.concat(['PUT'], arguments));
***REMOVED***

__internal__.modules.http.delete = function() ***REMOVED***
	return __internal__.modules.http.do.apply(this, _.concat(['DELETE'], arguments));
***REMOVED***
