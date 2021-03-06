// --- CKEditor ---
var editor = CKEDITOR.replace(
    'campaignTemplate', {
        filebrowserBrowseUrl: '/assets/filemanager/index.html?config=../../filemanager.config',
        extraPlugins:
        'codemirror',
        codemirror: {
            autoCloseBrackets: true,
            autoCloseTags: true,
            autoFormatOnStart: true,
            autoFormatOnUncomment: true,
            continueComments: true,
            enableCodeFolding: true,
            enableCodeFormatting: true,
            enableSearchTools: true,
            highlightMatches: true,
            indentWithTabs: false,
            lineNumbers: true,
            lineWrapping: true,
            mode: 'htmlmixed',
            matchBrackets: true,
            matchTags: true,
            showAutoCompleteButton: true,
            showCommentButton: true,
            showFormatButton: true,
            showSearchButton: true,
            showTrailingSpace: true,
            showUncommentButton: true,
            styleActiveLine: true,
            theme: 'default',
            useBeautifyOnStart: false
        },
        plugins:
        //'dialogui,' +
        'dialog,' +
        //'a11yhelp,' +
        //'dialogadvtab,' +
        'basicstyles,' +
        //'bidi,' +
        //'blockquote,' +
        'clipboard,' +
        'button,' +
        'panelbutton,' +
        'panel,' +
        'floatpanel,' +
        'colorbutton,' +
        'colordialog,' +
        'templates,' +
        'menu,' +
        'contextmenu,' +
        //'div,' +
        'resize,' +
        'toolbar,' +
        'elementspath,' +
        'enterkey,' +
        'entities,' +
        'popup,' +
        'filebrowser,' +
        'find,' +
        'fakeobjects,' +
        //'flash,' +
        'floatingspace,' +
        'listblock,' +
        'richcombo,' +
        'font,' +
        //'forms,' +
        'format,' +
        'horizontalrule,' +
        'htmlwriter,' +
        //'iframe,' +
        'wysiwygarea,' +
        'image,' +
        'indent,' +
        'indentblock,' +
        'indentlist,' +
        //'smiley,' +
        'justify,' +
        'menubutton,' +
        //'language,' +
        'link,' +
        'list,' +
        'liststyle,' +
        'magicline,' +
        //'maximize,' +
        //'pagebreak,' +
        'pastetext,' +
        'pastefromword,' +
        //'preview,' +
        //'print,' +
        'removeformat,' +
        'selectall,' +
        'showblocks,' +
        'showborders,' +
        'sourcearea,' +
        'specialchar,' +
        //'scayt,' +
        'stylescombo,' +
        'tab,' +
        'table,' +
        'tabletools,' +
        'undo,' +
        //'wsc,' +
        'docprops,',

        allowedContent:  true,
        removeFormatAttributes: '',
        height: 450,
        entities:false
    }
);

/*
editor.on( 'instanceReady', function() {
    editor.resize("100%", $("#template").height()-24)
} );
*/
// --- /CKEditor ---