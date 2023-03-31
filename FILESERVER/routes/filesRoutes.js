// Requiring express to handle routing
const express = require("express");
var path = require('path');
// The fileUpload npm package for handling
// file upload functionality
const fileUpload = require("express-fileupload");
  
// Creating app
var router = express.Router();
var cache_file_name;
  
// Passing fileUpload as a middleware
router.use(fileUpload());

router.get("/", function(req,res){
  
  res.sendFile(path.join(__dirname+"/../public/files.html"));

})
  
// For handling the upload request
router.post("/upload", function (req, res) {
  
  // When a file has been uploaded
  if (req.files && Object.keys(req.files).length !== 0) {
    
    // Uploaded path
    const uploadedFile = req.files.uploadFile;
  
    // Logging uploading file
    console.log(uploadedFile);
  
    // Upload path
    const uploadPath = __dirname
        + "/files/" + uploadedFile.name;
  
    // To save the file using mv() function
    uploadedFile.mv(uploadPath, function (err) {
      if (err) {
        console.log(err);
        res.send("Failed !!");
      } else res.send("Successfully Uploaded !!");
    });
  } else res.send("No file uploaded !!");
});
  

// To handle the download file request
router.post("/download/", function (req, res) {
    console.log(req.body.file_name);
    cache_file_name = req.body.file_name;

});

// To handle the download file request
router.get("/download/", function (req, res) {
   
  // The res.download() talking file path to be downloaded
  res.download(__dirname + "/files/"+cache_file_name, function (err) {
    if (err) {
      console.log(err);
    }
  });
});
  


module.exports = router;