var gulp = require('gulp')
var exec = require('child_process').exec;

function bower_install(cb){
  exec('bower install', function (err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    cb(err);
  });
}

function move_dist_files(cb) {

  console.log("Moving all files to templates folder");
  gulp.src("bower_components/jquery/dist/jquery.min.js")
      .pipe(gulp.dest('templates/js'));

  gulp.src("bower_components/bootstrap/dist/js/bootstrap.min.js")
      .pipe(gulp.dest('templates/js'));

  gulp.src("bower_components/bootstrap/dist/css/bootstrap.min.css")
      .pipe(gulp.dest('templates/css'));

  cb()
};

module.exports.default = gulp.series(bower_install, move_dist_files)
