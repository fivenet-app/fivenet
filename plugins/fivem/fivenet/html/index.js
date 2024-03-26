// General
let webUrl = 'https://example.com';
const animationSpeed = 325;

let tabletLoaded = false;

function closeAll() {
    $.post(`https://${GetParentResourceName()}/closeTablet`, JSON.stringify({ ok: true }));
    fadeOut('#tokenmgmt');
    fadeOut('#pwreset');
    closeTablet();
}

function fadeIn(id) {
    $(id).fadeIn(animationSpeed).css('display', 'block');
}

function fadeOut(id) {
    $(id).fadeOut(animationSpeed).css('display', 'none');
}

function copyToClipboard(text) {
    const el = document.createElement('textarea');
    el.value = text;
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
}

// Registration Token / Account Management
function closeTokenMgmt() {
    fadeOut('#tokenmgmt');
    $.post(`https://${GetParentResourceName()}/exit`, JSON.stringify({ ok: true }));
}

function setWebUrl(url) {
    webUrl = url;

    $('.fivenet_url').attr('href', url);
    $('.fivenet_url').text(url);
}

function openFiveNetWebsite(path) {
    window.invokeNative('openUrl', webUrl + (path === undefined ? '' : path));
}

// Tablet
function openTablet() {
    $('#tablet').animate({ top: '50%' }, animationSpeed).css('display', 'block');
}

function closeTablet() {
    $('#tablet').animate({ top: '150%' }, animationSpeed).css('display', 'block');
}

function foldTablet() {
    $('#tablet').toggleClass('tablet-flipped');
    if ($('#tablet').hasClass('tablet-flipped')) {
        $('.foldbutton').text('Zuklappen');
    } else {
        $('.foldbutton').text('Aufklappen');
    }
}

function refreshTablet() {
    $('#tablet-frame').attr('src', `${webUrl}?refreshApp=true&nui=${GetParentResourceName()}`);
}

function navigateTabletTo(route) {
    postMessageToTablet({ type: 'navigateTo', route });
}

function postMessageToTablet(data) {
    document.getElementById('tablet-frame').contentWindow.postMessage(data, webUrl);
}

// Message handler
$(document).ready(function () {
    window.addEventListener('message', (event) => {
        const item = event.data;

        if (item === undefined) {
            return;
        }

        switch (item.type) {
            // Registration Token / Account Management
            case 'token':
                setWebUrl(item.webUrl);

                // Hide Tablet and show Token UI
                closeTablet();

                fadeIn('#tokenmgmt');
                fadeOut('#pwreset');

                $('.token').text(item.data);
                $('.hint').html(
                    '<i class="fa-sharp fa-solid fa-circle-info"></i> Nutze diesen Token, um deinen Account zu registrieren oder dein Passwort zurückzusetzen.',
                );
                break;

            case 'username':
                setWebUrl(item.webUrl);

                // Hide Tablet and show Token UI
                closeTablet();

                fadeIn('#tokenmgmt');
                fadeIn('#pwreset');
                $('.token').text(item.data);
                $('.hint').html(
                    '<i class="fa-sharp fa-solid fa-circle-info"></i>Dein Account wurde mittels dieses Usernames registriert.',
                );
                break;

            // Tablet
            case 'openTablet':
                if (!tabletLoaded) {
                    setWebUrl(item.webUrl);

                    refreshTablet();
                    tabletLoaded = true;
                }

                // Hide Token UI
                fadeOut('#tokenmgmt');
                fadeOut('#pwreset');

                // Show Tablet UI
                openTablet();
                $(document).on('click.hideTabletClick', function (e) {
                    const container = $('#tablet');
                    if (!container.is(e.target) && container.has(e.target).length === 0) {
                        closeAll();
                    }
                });
                break;

            case 'closeTablet':
                // Hide Tablet UI and unbind click outside tablet
                $(document).unbind('click.hideTabletClick');
                closeTablet();
                break;

            case 'copyToClipboard':
                copyToClipboard(item.text);
                break;

            case 'fixTablet':
                console.info('Attempting to fix tablet...');

                fadeOut('#tokenmgmt');
                fadeOut('#pwreset');
                closeTablet('#tablet');

                // Navigate to clear site data and then back to FiveNet
                $('#tablet-frame').attr('src', `${item.webUrl}/api/clear-site-data`);

                setTimeout(() => {
                    console.info('Reloading FiveNet on tablet after fix ...');

                    refreshTablet();
                }, 3000);
                break;

            default:
                // If the tablet is loaded, forward any other messages to the tablet iframe
                if (tabletLoaded) {
                    postMessageToTablet(item);
                }
                break;
        }
    });

    $(window).on('keydown', function (event) {
        if (event.key === 'Escape' || event.key === 'F5') {
            closeAll();
        }
    });
});
