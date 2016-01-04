var gulp = require('gulp'),
    livereload = require('gulp-livereload'),
    watch = require('gulp-watch'),
    less = require('gulp-less'),
    connect = require('gulp-connect'),
    embedlr = require("gulp-embedlr");

var LIVEPORT = 35701;
var BUILD_DIR = './build'

gulp.task('less', function() {
	return gulp.src('./*.less')
	.pipe(less())
	.pipe(gulp.dest(BUILD_DIR))
	.pipe(livereload())
	;
});

gulp.task('inject_livereload', function() {
    gulp.src('./*.html')
    .pipe(embedlr({port: LIVEPORT}))
	.pipe(gulp.dest(BUILD_DIR))
    .pipe(livereload())
    ;
});

gulp.task('watch', function() {
    livereload.listen({port: LIVEPORT});
    gulp.watch('./*.less', ['less'])
    gulp.watch('./*.html', ['inject_livereload'])
});

gulp.task('server', function(){
    connect.server({
        root: BUILD_DIR,
        host: '127.0.0.1',
        port: 8030,
        livereload: LIVEPORT
    });
});

gulp.task('default', ['less', 'inject_livereload', 'watch', 'server']);
