
$(document).ready(function() 
{
    $('#menunav li').click(function(e) 
    { 
        $("#menunav li").removeClass("active");
        $("#menunav li").eq($(this).attr("id")).addClass('active');
    });
});

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
    