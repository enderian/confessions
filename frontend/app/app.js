import './app.scss'
import './tick.scss'

import 'bootstrap'

import './form'
import './index'
import $ from 'jquery'

$('#cookieAlert').on('closed.bs.alert', function () {
    document.cookie = 'cookie-alert-closed=1; max-age=31557600; path=/'
})

$(function () {
    if (!document.cookie.includes("cookie-alert-closed=1")) {
        $('#cookieAlert').show()
    }
})