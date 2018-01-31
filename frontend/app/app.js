import './app.scss'
import './tick.scss'
import 'bootstrap'
import $ from 'jquery'

$(function () {
    $('.custom-value').change(function () {
        const id = '#' + $(this).attr('name') + '-custom';
        if ($(this).val() === "custom") {
            $(id).show();
        } else {
            $(id).hide();
        }
    });

    $('#status_check_form').submit(function (event) {
        event.preventDefault();
        $.post('/secret', $('#status_check_form').serialize(), function(resp){
            $('#status_check_response').html(resp).show()
        })
    });

    document.deleteSubmit = function () {
        $.post('/secret', $('#status_check_delete').serialize(), function(resp){
            $('#status_check_response').html(resp).show()
        });
        return false
    };
});