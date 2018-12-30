# go-server


  This is a tcp server project coded with go. It doesn't include a client, if you want to work with it you need to use a client.  
It supports multiple clients at once.    


  -How to use it  
 After setting up a client, you need to send a tcp packet which contains reg,ID,PASSWORD to register and add your user to the JSON
file.  
If you want to skip this step and use the test accounts you can send log,ID,PASSWORD (You can find test accounts below). 
After logging in you can use handleCommands function to modify the users capabilities.

Server closes connection on these circumstances  
1-If user sends wrong password.  
2-If user sends wrong id.  
3-If user sends undefined string in packet before logging in.  


Test accounts;  
ID:  
test  
test2  
test3  
Password:test  

Example packet:"log,test,test" to log in.  
Default port is 7001.
