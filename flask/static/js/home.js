
$(document).ready(function () {
    $('#sidebarCollapse').on('click', function () {
        $('#sidebar').toggleClass('active');
        $(this).toggleClass('active');
    });
});

$(function() { $('textarea').froalaEditor() });     
$(document).ready(function(){
    $('[data-toggle="tooltip"]').tooltip(); 
});

function submitResult() {
    if ( confirm("Are you sure you wish to remove event?") == false ) {
       return false ;
    } else {
       return true ;
    }
 }