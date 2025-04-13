/*jslint browser: true*/ /*global  $*/ /*global  jQuery*/
$(document).ready(function(){
    const $navMenu = $("#mobileMenu");
    jQuery('.ntgl').click(function(){
        $navMenu.toggle("slide", { direction: "down" }, 600);
    });

    const updateSliderCounter = function (current, total){
        current = (
            current < 10 ? "0" + current : current);
        total = (
            total < 10 ? "0" + total : total);
        $(".slick-review-counter").text(current+"/"+total);

    };

    const $slickReview = $("#slick-review");

    //(ev Необходим для работы slick, но jsl
    //считает ошибкой то, что он не используется внутри)
    $slickReview.on("init", function (ev, slick){ //jslint-ignore-line
        updateSliderCounter(slick.slickCurrentSlide() + 1, slick.slideCount);
    });


    $slickReview.slick({
        dots: false,
        nextArrow: '',
        prevArrow: ''
    });
    
    //(ev Необходим для работы slick, но jsl
    //считает ошибкой то, что он не используется внутри)
    $slickReview.on('afterChange', function (ev, slick){ //jslint-ignore-line
        updateSliderCounter(slick.slickCurrentSlide()+1, slick.slideCount);
    });

    $(".review-arrow-prev").click(function (){
        $("#slick-review").slick("slickPrev");
    });
    $(".review-arrow-next").click(function (){
        $("#slick-review").slick("slickNext");
    });
    $('.panel-btn').click(function (){
        const t = $(this).parents(".panel");
        t.toggleClass("panel_open");
        t.toggleClass("panel_close");
        t.children(".panel-body").slideToggle(400);
    });
    $('#slick-customers-first').slick({
        autoplay: true,
        autoplaySpeed: 2000,
        infinite: true,
        nextArrow: ``,
        prevArrow: ``,
        responsive: [
            {
                breakpoint: 1024,
                settings: {
                    infinite: true,
                    slidesToScroll: 1,
                    slidesToShow: 5
                }
            },
            {
                breakpoint: 600,
                settings: {
                    slidesToScroll: 2,
                    slidesToShow: 4
                }
            },
            {
                breakpoint: 480,
                settings: {
                    slidesToScroll: 1,
                    slidesToShow: 3
                }
            }
            // You can unslick at a given breakpoint now by adding:
            // settings: "unslick"
            // instead of a settings object
        ],
        slidesToScroll: 1,
        slidesToShow: 6,
        speed: 600
    });
    $('#slick-customers-second').slick({
        autoplay: true,
        autoplaySpeed: 3250,
        infinite: true,
        nextArrow: ``,
        prevArrow: ``,
        responsive: [
            {
                breakpoint: 1024,
                settings: {
                    infinite: true,
                    slidesToScroll: 1,
                    slidesToShow: 5
                }
            },
            {
                breakpoint: 600,
                settings: {
                    slidesToScroll: 2,
                    slidesToShow: 4
                }
            },
            {
                breakpoint: 480,
                settings: {
                    slidesToScroll: 1,
                    slidesToShow: 3
                }
            }
            // You can unslick at a given breakpoint now by adding:
            // settings: "unslick"
            // instead of a settings object
        ],
        slidesToScroll: 1,
        slidesToShow: 6,
        speed: 600
    });

    
});