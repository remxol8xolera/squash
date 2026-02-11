package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	// "log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type LLMResponse struct {
	LLMResponse string `json:"llm_response"`
}
type error_response_json struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Username   string `json:"username"`
}
type error_response_json_for_django_backend struct {
	Error_message        string `json:"error_message"`
	Message_for_the_user string `json:"message_for_the_user"`
	StatusCode           int    `json:"status_code"`
	Username             string `json:"username"`
}
type response_json_for_django_backend struct {
	Error_message                    string `json:"error_message"`
	Link_for_the_current_site        string `json:"link_for_the_current_site"`
	Message_for_the_user             string `json:"message_for_the_user"`
	StatusCode                       int    `json:"status_code"`
	Username                         string `json:"username"`
}
type error_response_json_for_django_backend_with_an_array_as_value struct {
	Error_message        string `json:"error_message"`
	Message_for_the_user string `json:"message_for_the_user"`
	StatusCode           int    `json:"status_code"`
	Username             string `json:"username"`
	Values              []string `json:"values"`
}
type json_error_response_query_not_present struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("got /hello request\n")

	io.WriteString(w, "Hello, World from me !\n")
}
func get_all_the_projects_of_the_user(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { //-------------- change it -------------------
		print("Oh  my god")
		return_json_error(w, http.StatusMethodNotAllowed, error_response_json{
			Error:      "method not allowed",
			Message:    "only get method is allowed on this route",
			StatusCode: http.StatusMethodNotAllowed,
			Username:   "",
		})
		return
	}
	if validate_url_params_if_not_present_write_bad_request(r, w, "userName") {
		// if true meaning bad request return
		return
	}
	query := r.URL.Query()
	userName := query.Get("userName")
	dir, err := os.Open("src/routes/"+userName)
	if err!= nil {
		if err.Error() == "open src/routes/"+userName+": no such file or directory"{
			// --------------------- wait ---------------
			// if the user name is not found , that meand that the backend send such response or that usr does not exist , or the dir for the user is not created yet
			//  check if the user name is ../ they can see the files in the dir 
			// --------------------- wait ---------------
			return_json_error(w, http.StatusBadRequest, error_response_json_for_django_backend{
				Error_message:     		   " dir with that username not found  -->> "+userName+"<<-- dir.",
				Message_for_the_user: 	   "Oops! Your name was not found on the server ",
				StatusCode: 			   http.StatusBadRequest, // http.StatusNotFound (404) <<- is the good one here   
				Username: 				   userName,
			}) 
			return
		}


	
		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
			Error_message:     		   " can't open the userName  dir for -->> "+userName+"<<-- dir.\n error Backend got-->>  "+err.Error(),
			Message_for_the_user: 	   "Oops! An error occurred on our side , can't get the names of your projects",
			StatusCode: 			   http.StatusInternalServerError, // http.StatusNotFound (404) <<- is the good one here   
			Username: 				   userName,
		}) 
		return
	}
	defer dir.Close()
	println(dir)
	things_in_dir,error_form_reading_dir := dir.ReadDir(-1)
	if error_form_reading_dir != nil {
		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend_with_an_array_as_value{
			Error_message:     		   "can't open the userName  dir for -->> "+userName+"<<-- dir.\n error Backend got-->>  "+err.Error(),
			Message_for_the_user: 	   "Oops! An error occurred on our side , can't get the names of your projects",
			StatusCode: 			   http.StatusInternalServerError, // http.StatusNotFound (404) <<- is the good one here   
			Username: 				   userName,
			Values: 					[]string{},
		}) 
		return
	}
	var directories []string
    for _, entry := range things_in_dir {
        if entry.IsDir() {
			println("name of the dir ",entry.Name())
            directories = append(directories, entry.Name())
        }
    }
	
	return_json_error(w, http.StatusOK, error_response_json_for_django_backend_with_an_array_as_value{
		Error_message:     "" ,
		Message_for_the_user: "",
		StatusCode:        http.StatusOK,
		Username:          userName,
		Values:       directories,
	})


}


func delete_a_project(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete { //-------------- change it -------------------
		print("Oh  my god")
		return_json_error(w, http.StatusMethodNotAllowed, error_response_json{
			Error:      "method not allowed",
			Message:    "only post method is allowed on this route",
			StatusCode: http.StatusMethodNotAllowed,
			Username:   "",
		})
		return
	}
	if validate_url_params_if_not_present_write_bad_request(r, w, "userName") {
		// if true meaning bad request return
		return
	}
	if validate_url_params_if_not_present_write_bad_request(r, w, "project_name") {
		// if true meaning bad request return
		return
	}
	query := r.URL.Query()
	userName := query.Get("userName")
	var project_name  = query.Get("project_name")
	if project_name == "temp"{
		return_json_error(w, http.StatusBadRequest, error_response_json_for_django_backend{
			Error_message:     		   " can't delete the project -->> "+project_name+"<<-- dir.",
			Message_for_the_user: 	   "Oops! can't delete the temp, if you want to change it generate a new one ",
			StatusCode: 			   http.StatusBadRequest, // http.StatusNotFound (404) <<- is the good one here   
			Username: 				   userName,
		}) 
		return
	}
	dir_entries , error :=os.ReadDir("src/routes/"+userName+"/"+project_name)
	if error!= nil {
		
		if error.Error() == "open src/routes/"+userName+"/"+project_name+": no such file or directory"{
			return_json_error(w, http.StatusBadRequest, error_response_json_for_django_backend{
				Error_message:     		   " can't delete the project -->> "+project_name+"<<-- dir.",
				Message_for_the_user: 	   "Oops! A project with that name was not found ",
				StatusCode: 			   http.StatusBadRequest, // http.StatusNotFound (404) <<- is the good one here   
				Username: 				   userName,
			}) 
			return
		}
		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
			Error_message:     		   " can't delete the project -->> "+project_name+"<<-- dir.  ",
			Message_for_the_user: 	   "Oops! A project with that name was not found ",
			StatusCode: 			   http.StatusInternalServerError, // http.StatusNotFound (404) <<- is the good one here   
			Username: 				   userName,
		}) 
		return
	}
	println(dir_entries)

	err:= os.RemoveAll("src/routes/"+userName+"/"+project_name)
	if err != nil {
		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
			Error_message:     		   " can't delete the project -->> "+project_name+"<<-- dir.  ",
			Message_for_the_user: 	   "Oops! an error occured on our side while deleting your project  ",
			StatusCode: 			   http.StatusInternalServerError,
			Username: 				   userName,
		}) 
		return
		
	}
	return_json_error(w, http.StatusOK, error_response_json_for_django_backend{
		Error_message:     		   " Successfully  deleted the project_name -->> "+project_name+"<<--  dir.  ",
		Message_for_the_user: 	   "Successfully deleted your project",
		StatusCode: 			   http.StatusOK,
		Username: 				   userName,
	}) 

}


func host_the_temp_one_in_a_production_site(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //-------------- change it -------------------
		print("Oh  my god")
		return_json_error(w, http.StatusMethodNotAllowed, error_response_json{
			Error:      "method not allowed",
			Message:    "only post method is allowed on this route",
			StatusCode: http.StatusMethodNotAllowed,
			Username:   "",
		})
		return
	}
	// ------------modify the function to take a pointer to the value and add it there  -- so i don't have to keep the 
	// ------------searching the url again  and again 
	if validate_url_params_if_not_present_write_bad_request(r, w, "userName") {
		// if true meaning bad request return
		return
	}
	// println("\n query := r.URL.Query()",r.URL.Query().Get("userName"))
	if validate_url_params_if_not_present_write_bad_request(r, w, "project_name") {
		// if true meaning bad request return
		return
	}
	query := r.URL.Query()
	userName := query.Get("userName")
	var project_name  = query.Get("project_name")
	//  checked for the userName and project_name now just copy the file from  the temp and make a project name dir and add to it 
	err := create_dir("src/routes/"+userName,project_name)
	if err != nil{

		// ------------------------ wait ---------------------------------
		// 	
		// we can also write to it , or just ignore it -- i think if user want to go for the same they should deleate and create a new one 
		// that can even be with the same one wit 
		// 
		// 2.> what waht if thetemp dir is empty (it is not , i meant what if user tryied to push trial page that i created to the project )
		//  well it will not happen as I give them the their first website as a temlate 
		// 
		//  3.> what if the user name and the project name is wrong and does not exist 
		// ------------------------ wait ---------------------------------
		
		if "mkdir src/routes/"+userName+"/"+project_name+": file exists" == err.Error(){
			return_json_error(w, http.StatusBadRequest, error_response_json_for_django_backend{
				Error_message:     		   "name chosen by the user--"+userName+" is same as "+project_name,
				Message_for_the_user: 	   "Project with that name already exists, please chose another name or delete the project with that name first ",
				StatusCode: 			   http.StatusBadRequest,
				Username: 				   userName,
			}) 
			return

		}else{

			return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
				Error_message:     		   "unable to create the project dir. ",
				Message_for_the_user: 	   "Unable to host your website , please try again , if that does not work try a different project name ",
				StatusCode: 			   http.StatusInternalServerError,
				Username: 				   userName,
			})
			return
		}
	}
	// reading the file form the temp dir 
	file_in_the_temp_dir , err := os.ReadFile("src/routes/"+userName+"/temp/+page.svelte")
	if err != nil {
		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
			Error_message:     		   "unable to open the temp dir of "+userName,
			Message_for_the_user: 	   "Oops! an error occured on our side , don't worry try again ",
			StatusCode: 			   http.StatusInternalServerError,
			Username: 				   userName,
		})
		return
	} 
	//  -------------- what if the file already exists(project name ) and what if the user is not created  -- auth is done on 
	//  -------------- the django backend so trust  it and create it again  -->  check if it exist it not create it if it does prompt 
	// --------------- the user to provide new pjectname
	// -----------   ---------
	// ----- err.Error() is -> mkdir src/routes/monish/erhb: file exists
	// -----------probally should also delete the project dir if unsuccessful giving user what they want as they should be able to  
	file, err := os.Create(filepath.Join("src/routes/"+userName+"/"+project_name,"+page.svelte" ))
	if err != nil {
		os.RemoveAll("src/routes/"+userName+"/"+project_name) // --------------------||==>>this one could error so keep that in mind
		println("in here , i am  in jail irong lung ")
		//  if i can not create a file in this dir ; if thename chosen by the user is same (as project name ) I am already returning an
		//  error so that maeans if the request is comming here it is new -> we should create a file (new one ) and that means 
		//  here we should delete this dir
		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
			Error_message:     		   "unable to open the temp dir of "+userName,
			Message_for_the_user: 	   "Oops! an error occured on our side , don't worry try again ",
			StatusCode: 			   http.StatusInternalServerError,
			Username: 				   userName,
		})
		return
	}

	_, error_122:= file.WriteString(string(file_in_the_temp_dir))	
	if error_122!= nil{

		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
			Error_message:     		   "unable to write to the production dir  of the app   "+project_name+" for ->"+userName,
			Message_for_the_user: 	   "Oops! an error occured on our side , don't worry try publishing again or create a new one  ",
			StatusCode: 			   http.StatusInternalServerError,
			Username: 				   userName,
		})
	}
	println("\n\nos env -->>",os.Getenv("SVELTE_URL_WITH_SLASH"))
	
	return_json_error(w, http.StatusOK, response_json_for_django_backend{
		Error_message:     		   "successfully made the project  "+project_name+" for "+userName,
		Link_for_the_current_site: os.Getenv("SVELTE_URL_WITH_SLASH")+userName+"/"+project_name,
		Message_for_the_user: 	   "Yay! your site is created successfully and is live ",
		StatusCode: 			   http.StatusOK,
		Username: 				   userName,
	})


}

func llm_response_write_it_in_temp_dir(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userName := query.Get("userName")
	if r.Method != http.MethodPost { 
		print("Oh  my god")
		return_json_error(w, http.StatusMethodNotAllowed, error_response_json{
			Error:      "method not allowed",
			Message:    "only post method is allowed on this route",
			StatusCode: http.StatusMethodNotAllowed,
			Username:   "",
		})
		return
	}
	if validate_url_params_if_not_present_write_bad_request(r, w, "userName") {
		// making it look for the userName param in the url and not the userName
		// if true meaning bad request return
		return
	}
	var llmResponseData LLMResponse
	if get_json_field_out_of_body_and_write_error_on_response(w, r, &llmResponseData) {
		print("i think llmResponse is -->", llmResponseData.LLMResponse)
		print("\n ======= from json field servhing in the body function")
		return
	}
	// --------------- break this func down to 1.> decoding json  and 2.> valaditing decoded json ------------
	// --------------- that way we will be able to take the decoded json and write it to the file
	// --------------------------------------------------------
	// ------------------instead I will decode it once again to , and this time there will be no  error , as it is already been solved  previously
	// ----------------------------or---------------------------------------
	// ---------------------------just pass a pointer

	// --------2nd step , write it to a file-->> could also create it if it does not exist (username dir. )
	// os_file , err := os.OpenFile("src/routes/"+userName+"/temp/+page.svelte", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	os_file, err := os.Create(filepath.Join("src/routes/"+userName+"/temp", "+page.svelte"))
	if err != nil {
		println("\n -- this functionn above can return error , whereas we should try to make the dirs. and file inside it here , make a integration test for it ---")
		// return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
		// 	Error_message:        "can't open the file for you ",
		// 	Message_for_the_user: "Sorry we are having trouble keeping our side with us , please try again ",
		// 	StatusCode:           http.StatusInternalServerError,
		// 	Username:             userName,
		// })
		
		//  well if a user is sent by the djnago and the user does not exist here  then that means something is wrong here  , may be svelte
		// was down  , so i think we should create it here right not
		// panic(err)

		// userName dir
		err := create_dir("src/routes/", userName)
		println("\nabout to create a username dir")
		if err != nil {
			println("\nabout to create a username dir")
			return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
				Error_message:        "can't open the file for you ",
				Message_for_the_user: "Sorry we are having trouble keeping our side with us , please try to login again  ",
				StatusCode:           http.StatusInternalServerError,
				Username:             userName,
			})
			return
		}
		// temp dir
		println("\n about to create the  temp dir")
		erro := create_dir("src/routes/"+userName, "temp")
		if erro != nil {
			println("\nin the error of  temp dir",erro.Error())
			return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
				Error_message:        "can't open the file for you ",
				Message_for_the_user: "Sorry we are having trouble keeping our side with us , please try to login again , that should probally fix it   ",
				StatusCode:           http.StatusInternalServerError,
				Username:             userName,
			})
			return
		}
		println("\n about to create a +pages.svelte file ")
		error_l := only_create_file("+page.svelte", "src/routes/"+userName+"/temp")
		if error_l != nil {
			println("\nin the error of creating a svelte file", error_l.Error())
			return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
				Error_message:        "can't open the file for you ",
				Message_for_the_user: "Sorry we are having trouble keeping our side with us , please try to login again , that should probally fix it   ",
				StatusCode:           http.StatusInternalServerError,
				Username:             userName,
			})
			return
		}
		
		//  now if the func has not returned (or err != nil in any one ), that means we were able to create the 
		// 	dir once again and now lets wtite to it 
		os_file2, err2 := os.OpenFile("src/routes/"+userName+"/temp/+page.svelte", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err2 != nil {
			// this time just return the json error
			return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
				Error_message:        "failed to find  the user dir and even tried to create it once ",
				Message_for_the_user: "Sorry we are having trouble keeping our side with us , please try to login again , that should probally fix it   ",
				StatusCode:           http.StatusInternalServerError,
				Username:             userName,
			})
			return
		}
		defer os_file2.Close()
		println("\nabout to write this to a file --> ", llmResponseData.LLMResponse)
		os_file2.WriteString(llmResponseData.LLMResponse)
		// once the file is done writing to it return
		return_json_error(w, http.StatusCreated, error_response_json_for_django_backend{
			Error_message:        "",
			Message_for_the_user: "Successfully create your website ",
			StatusCode:           http.StatusCreated,
			Username:             userName,
		})
		return
	}

	defer os_file.Close()
	println("\n about to write to the file from the first func ")
	os_file.WriteString(llmResponseData.LLMResponse)
	return_json_error(w, http.StatusCreated, response_json_for_django_backend{
		Error_message:        "",
		Message_for_the_user: "Successfully create your website ",
		Link_for_the_current_site: os.Getenv("SVELTE_URL_WITH_SLASH")+userName+"/temp",
		StatusCode:           http.StatusCreated,
		Username:             userName,
	})

}


func create_temp_and_name_dir_for_user(w http.ResponseWriter, r *http.Request) {
	println("\n \n ------- form the temp dir for the new user func --------\n\n")
	if r.Method != http.MethodPost { 
		print("Oh  my god")
		return_json_error(w, http.StatusMethodNotAllowed, error_response_json{
			Error:      "method not allowed",
			Message:    "only post method is allowed on this route",
			StatusCode: http.StatusMethodNotAllowed,
			Username:   "",
		})
		return
	}

	query := r.URL.Query()
	userName := query.Get("userName")
	if userName == "" {
		// http.Error(w, "userName not provided ", http.StatusBadRequest)
		return_json_error(w, http.StatusNotAcceptable, error_response_json_for_django_backend{
			Error_message:        "userName not provided",
			Message_for_the_user: "userName not provided, if the error continue  you should try logging in again  ",
			StatusCode:           http.StatusNotAcceptable,
			Username:             userName,
		})

		// ------------------------ potential error ------------------------------------
		// >> (<-means done) i should tell my backend this error is for it , meaning  it has not provided the username , user should not be alerted
		// ------------------------ potential error ------------------------------------
		return
	}
	// -----------------

	error := create_dir("src/routes", userName)
	if error != nil {
		var err_if_dir_is_already_there = "mkdir src/routes/" + userName + ": file exists"
		fmt.Printf("here %s \n\n", userName)
		if error.Error() != err_if_dir_is_already_there {
			return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
				Error_message:        "failed to create the username dir  ",
				Message_for_the_user: "userName not provided, if the error continue  you should try logging in again  ",
				StatusCode:           http.StatusInternalServerError,
				Username:             userName,
			})  
			return
		}
		// well if the user name is created keep looking in it to check for the other dir (just return 200 or do  not at all-- if it already exista )
	}

	// creating the temp dir
	error_from_temp_dir := create_dir("src/routes/"+userName, "temp")
	if error_from_temp_dir != nil {
		print(error_from_temp_dir.Error(), "\n\n")
		// var err_if_dir_is_already_there  = "mkdir src/routes/"+userName+"/temp: file exists"
		var err_if_dir_is_already_there = "mkdir src/routes/" + userName + ": file exists" // ---bro i don't get it the print statement shows/tells
		// why does this without /temp works , idk
		if error.Error() != err_if_dir_is_already_there {
			print("in the error which  is not about same --temp --dir ")
			return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
				Error_message:        "failed to create the temp dir for the "+userName,
				Message_for_the_user: "Oops! an error occured on our side while creating you account your account , Loggin again should probally solve it  ",
				StatusCode:           http.StatusInternalServerError,
				Username:             userName,
			})
			return
		}
	}
	// creating the file in temp dit
	error_by_creating__first_file_in_temp := only_create_file("+page.svelte", "src/routes/"+userName+"/temp")
	if error_by_creating__first_file_in_temp != nil {
		http.Error(w, error_by_creating__first_file_in_temp.Error()+"\n got the error creating the temp dir  ", http.StatusInternalServerError)
		return_json_error(w, http.StatusInternalServerError, error_response_json_for_django_backend{
			Error_message:        "failed to create the +pages.svelte file in temp dir for  the "+userName,
			Message_for_the_user: "Oops! an error occured on our side while creating you account your account , Loggin again should probally solve it",
			StatusCode:           http.StatusInternalServerError,
			Username:             userName,
		})
		return
	}
	// http.Response("here is your file "+names_of_the_file,200)
	return_json_error(w, http.StatusCreated, error_response_json_for_django_backend{
		Error_message:        "successfully created the user ",
		Message_for_the_user: "Successfully made your account , your temp website is live ",
		StatusCode:           http.StatusCreated,
		Username:             userName,
	})
}
func getHello2(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
  fmt.Printf("\nresponse from the go server -->>",w)
	return_json_error(w, http.StatusCreated, error_response_json_for_django_backend{
		Error_message:        "successfully created the user ",
		Message_for_the_user: "Successfully amde your accunt , your temp website is live ",
		StatusCode:           http.StatusCreated,
    Username:             "userName",
	})
	io.WriteString(w, "Hello, HTTP!\n")
}


func main() {
	println("i reached down ")
erro := godotenv.Load()
if erro != nil {
	  println("i reached down----- ")
    // log.Fatal("Error loading .env file")
  }
  println("trying the env in the main func", os.Getenv("SVELTE_URL_WITH_SLASH"))

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/create_temp_and_name_dir_for_user", create_temp_and_name_dir_for_user) // done in django 
	http.HandleFunc("/llm_response_write_it_in_temp_dir", llm_response_write_it_in_temp_dir) // done in django 
	http.HandleFunc("/host_the_temp_one_in_a_production_site", host_the_temp_one_in_a_production_site) // done in django 
	http.HandleFunc("/delete_a_project", delete_a_project)  // done in django 
	http.HandleFunc("/get_all_the_projects_of_the_user", get_all_the_projects_of_the_user) 
  http.HandleFunc("/store_llm_response_in_trial_dir",getHello2)

	fmt.Printf("\n\n  ----------- go server listening on port http://localhost:4696   -------------\n\n")
	err := http.ListenAndServe(":4696", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)

	}
	// erro := godotenv.Load()
    // if erro != nil {
    //     // log.Fatal("Error loading .env file")
	// 	panic("Error Loading the .env file")
    // }


}

// -------------------Helper function ----------------------------------
func get_json_field_out_of_body_and_write_error_on_response(w http.ResponseWriter, r *http.Request, json_response_variable interface{}) bool {
	print(r.Body, "\n body form the json ")

	err := json.NewDecoder(r.Body).Decode(json_response_variable)
	println("\n json response from ===", json_response_variable)

	if err != nil {
		// Handle JSON decoding error 
		return_json_error(w, http.StatusBadRequest, json_error_response_query_not_present{
			Error:      "Invalid JSON input",
			Message:    "llm_response not provided in the json in request body",
			StatusCode: http.StatusBadRequest,
		})
		return true
	}

	// Check if the "llm_response" field is present
	llmResponse, ok := json_response_variable.(*LLMResponse)
	if !ok || llmResponse.LLMResponse == "" {
		// Handle missing "llm_response" field
		return_json_error(w, http.StatusBadRequest, json_error_response_query_not_present{
			Error:      "missing or empty field",
			Message:    "The 'llm_response' field is required",
			StatusCode: http.StatusBadRequest,
		})
		return true
	}

	return false
}

func validate_url_params_if_not_present_write_bad_request(url_query *http.Request, w http.ResponseWriter, check_for string) bool {
	// true means bad request and false means not
	query := url_query.URL.Query()
	checked_string := query.Get(check_for)
	if checked_string == "" {
		return_json_error(w, http.StatusBadRequest, json_error_response_query_not_present{
			Error:      "Bad request ",
			Message:    check_for + " was not found in you url",
			StatusCode: http.StatusBadRequest,
		})
		return true
	} else {
		return false
	}
}

func return_json_error(w http.ResponseWriter, http_status_error int, error_response_json any) error {
	w.WriteHeader(http_status_error)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(
		// error_response_json{
		// 	Error: "method not allowed",
		// 	Message: "only post method is allowed on this route",
		// 	StatusCode: http.StatusMethodNotAllowed,
		// 	Username: userName,
		// }
		error_response_json,
	)
}

func create_dir(path string, name string) error {
	// println("\n path from the create dir func -->",path, " ==name ",name)
	err := os.Mkdir(path+"/"+name, os.ModePerm)
	if err != nil {
		// handle error
		return err
	}
	return nil
}

func only_create_file(name_of_the_file string, path string) error {
	// if the file contains the same content keep it there as writing to it or not doing it is same in both cases (here chose
	//  writing to make sure to retain the default state)
	file, err := os.Create(filepath.Join(path, name_of_the_file))
	if err != nil {
		return err
	}
	svelte_component := `
	
<main>
<div class="container">
  <h1>
	<span class="text-animation text">Start by describing us your website</span>
  </h1>
</div>
</main>

<style>
main {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh; /* This will make the main element take up the entire page height */
  background-color: #121212;
  color: #f5f5f5;
  font-family: sans-serif;
}

.container {
  text-align: center;
}

h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  position: relative;
}
.text{
 color: rgb(255, 2, 52);   
}

.text-animation {
  background: linear-gradient(
	90deg,
	rgb(255, 2, 120) 50%,
	rgb(1, 247, 165) 0%,
	rgb(1, 99, 247) 100%
  );
  background-size: 200% 200%;
  background-position: 100% 0;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: text-animation 3s ease 1 forwards;
}

@keyframes text-animation {
  0% {
	background-position: 100% 0;
  }
  50% {
	background-position: 0 0;
  }
  100% {
	background-position: -100% 0;
  }
}

</style>

	`
	file.WriteString(svelte_component)
	defer file.Close() // as i just want to create a file
	return nil
}

// -------------------Helper function ----------------------------------
