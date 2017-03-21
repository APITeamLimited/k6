//
// This is an advanced k6 script sample that simulates users
// logging into an e-commerce site and purchasing things there.
//

import ***REMOVED*** group, check, sleep, fail ***REMOVED*** from "k6";
import http from "k6/http";

export let options = ***REMOVED***
	maxRedirects: 10,
***REMOVED***

// Base URL where the site is located
const baseURL = "https://dev-li-david.pantheonsite.io";

// User think time in between page loads etc. (change this to 0 when debugging)
const thinkTime = 30;

// List of login usernames and passwords
const credentials = [
	***REMOVED*** username: "testuser1", password: "testuser1" ***REMOVED***,
];

// Main function that gets executed by VUs repeatedly
export default function() ***REMOVED***
	// Load front/home page on site
	group("front page", doFrontPage);
	sleep(thinkTime);

	// Go to user login page and log in
	group("login page", doLogin);
	sleep(thinkTime);

	// Visit a random selection of available product category pages, and 
	// randomly add products from each category to our cart
	for (name in categories) ***REMOVED***
		if (Math.random() <= categories[name].chance) ***REMOVED***
			group(name, function() ***REMOVED*** doCategory(categories[name]); ***REMOVED***);
			sleep(thinkTime);
		***REMOVED***
	***REMOVED***

	// Check out our cart, perform payment
	group("checkout", doCheckout);
	sleep(thinkTime);

	// Log out the user
	group("logout", doLogout);
	sleep(thinkTime);
***REMOVED***

// "categories" is an object containing information about the various 
// product categories and the products available in each category:
//
// categories[name].url      = the URL for the product category page.
//
// categories[name].chance   = the likelihood that the product category 
//                             gets visited by the VU on each VU iteration.
//
// categories[name].title    = the HTML <title> text to look for, to determine 
//                             that the category page was loaded correctly.
//
// categories[name].products = an object containing info about the different
//                             products available under that category. See
//                             below for a description of the product object.
//
//
// products[name].url    = the URL for the product page.
//
// products[name].chance = the likelihood that a) the product page gets 
//                         visited when its parent product category page
//                         gets visited, and b) the product gets added to
//                         the VUs shopping cart when the product page is
//                         visited.
//              
// products[name].title  = the HTML <title> text to look for, when trying
//                         to determine if the product page was loaded 
//                         correctly.
//
const categories = ***REMOVED***
	"To Carry": ***REMOVED***
		url: `$***REMOVED***baseURL***REMOVED***/collection/carry`,
		title: "To carry | David li commerce-test",
		chance: 0.5,
		products: ***REMOVED***
			"Laptop bag": ***REMOVED***
				url: `$***REMOVED***baseURL***REMOVED***/bags-cases/commerce-guys-laptop-bag`,
				title: "Commerce Guys Laptop Bag | David li commerce-test",
				chance: 0.25,
			***REMOVED***,
			"Drupal Bag": ***REMOVED***
				url: `$***REMOVED***baseURL***REMOVED***/bags-cases/drupal-commerce-messenger-bag`,
				title: "Drupal Commerce Messenger Bag | David li commerce-test",
				chance: 0.25,
			***REMOVED***,
		***REMOVED***
	***REMOVED***,
	"To Drink With": ***REMOVED***
		url: `$***REMOVED***baseURL***REMOVED***/collection/drink`,
		title: "To drink with | David li commerce-test",
		chance: 0.5,
		products: ***REMOVED***
			"Drupal Commerce to Wake You Up": ***REMOVED***
				url: `$***REMOVED***baseURL***REMOVED***/drinks/drupal-commerce-wake-you`,
				title: "Drupal Commerce to Wake You Up | David li commerce-test",
				chance: 0.25,
			***REMOVED***,
			"The Guy Mug": ***REMOVED***
				url: `$***REMOVED***baseURL***REMOVED***/drinks/guy-mug`,
				title: "\"The Guy\" Mug  | David li commerce-test",
				chance: 0.25,
			***REMOVED***,
		***REMOVED***
	***REMOVED***
***REMOVED***;


// Request the front/home page, plus static resources
function doFrontPage() ***REMOVED***
	// Load the front page.
	check(http.get(baseURL + "/"), ***REMOVED***
		"title is correct": (res) => res.html("title").text() == "Welcome to David li commerce-test | David li commerce-test",
	***REMOVED***);

	// Load static assets.
	group("static resources", function() ***REMOVED*** http.batch(staticAssets); ***REMOVED***);
***REMOVED***

// Request the login page and perform a user login
function doLogin() ***REMOVED***
	// Request the login page.
	let res = http.get(baseURL + "/user/login");
	check(res, ***REMOVED***
		"title is correct": (res) => res.html("title").text() == "User account | David li commerce-test",
	***REMOVED***);

	// TODO: Add attr() to k6/html!
	// Extract hidden input fields.
	let formBuildID = res.body.match('name="form_build_id" value="(.*)"')[1];

	group("login", function() ***REMOVED***
		// Pick a random set of credentials.
		let creds = credentials[Math.floor(Math.random()*credentials.length)];

		// Create login request.
		let formdata = ***REMOVED***
			name: creds.username,
			pass: creds.password,
			form_build_id: formBuildID,
			form_id: "user_login",
			op: "Log in",
		***REMOVED***;
		let headers = ***REMOVED*** "Content-Type": "application/x-www-form-urlencoded" ***REMOVED***;
		// Send login request
		let res = http.post(baseURL + "/user/login", formdata, ***REMOVED*** headers: headers ***REMOVED***);
		// Verify that we ended up on the user page
		check(res, ***REMOVED***
			"login succeeded": (res) => res.url == `$***REMOVED***baseURL***REMOVED***/users/$***REMOVED***creds.username***REMOVED***`,
		***REMOVED***) || fail("login failed");
	***REMOVED***);
***REMOVED***

// Load a product category page, then potentially load some product pages
function doCategory(category) ***REMOVED***
	check(http.get(category.url), ***REMOVED***
		"title is correct": (res) => res.html("title").text() == category.title,
	***REMOVED***);

	for (prodName in category.products) ***REMOVED***
		if (Math.random() <= category.products[prodName].chance) ***REMOVED***
			group(prodName, function() ***REMOVED*** doProductPage(category.products[prodName]) ***REMOVED***);
		***REMOVED***
	***REMOVED***
***REMOVED***

// Load product page and potentially add product to cart
function doProductPage(product) ***REMOVED***
	let res = http.get(product.url);
	check(res, ***REMOVED***
		"title is correct": (res) => res.html("title").text() == product.title,
	***REMOVED***);
	if (Math.random() <= product.chance) ***REMOVED***
		let formBuildID = res.body.match('name="form_build_id" value="(.*)"')[1];
		let formID = res.body.match('name="form_id" value="(.*)"')[1];
		let formToken = res.body.match('name="form_token" value="(.*)"')[1];
		let productID = res.body.match('name="product_id" value="(.*)"')[1];
		group("add to cart", function() ***REMOVED*** 
			addProductToCart(product.url, productID, formID, formBuildID, formToken)
		***REMOVED***);
	***REMOVED***
***REMOVED***

// Add a product to our shopping cart
function addProductToCart(url, productID, formID, formBuildID, formToken) ***REMOVED***
	let formdata = ***REMOVED***
		product_id: productID,
		form_id: formID,
		form_build_id: formBuildID,
		form_token: formToken,
		quantity: 1,
		op: "Add to cart",
	***REMOVED***;
	let headers = ***REMOVED*** "Content-Type": "application/x-www-form-urlencoded" ***REMOVED***;
	let res = http.post(url, formdata, ***REMOVED*** headers: headers ***REMOVED***);
	// verify add to cart succeeded
	check(res, ***REMOVED***
		"add to cart succeeded": (res) => res.body.includes('Item successfully added to your cart')
	***REMOVED***) || fail("add to cart failed");
***REMOVED***

// Perform multi-step (multi-page) checkout
function doCheckout() ***REMOVED***
	var res;

	group("Checkout 1: review cart", function() ***REMOVED***
		// Checkout step 1: Review cart
		res = http.get(baseURL + "/cart");
		check(res, ***REMOVED***
			"title is correct": (res) => res.html("title").text() === 'Shopping cart | David li commerce-test'
		***REMOVED***);
	***REMOVED***);

	// Did we add any products to our cart?  If not, skip the rest of the checkout process
	if (res.body.includes("Your shopping cart is empty.")) ***REMOVED***
		return;
	***REMOVED***

	group("Checkout 2: submit cart", function() ***REMOVED***
		let formID = res.body.match('name="form_id" value="(.*)"')[1];
		let formToken = res.body.match('name="form_token" value="(.*)"')[1];
		let formBuildID = res.body.match('name="form_build_id" value="(.*)"')[1];
		let formdata = ***REMOVED***
			"form_build_id": formBuildID,
			"form_token": formToken,
			"form_id": formID,
			"op": "Checkout"
		***REMOVED***;
		let headers = ***REMOVED*** "Content-Type": "application/x-www-form-urlencoded" ***REMOVED***;
		res = http.post(baseURL + "/cart", formdata, ***REMOVED*** headers: headers ***REMOVED***);
		check(res, ***REMOVED***
			"cart submit succeeded": (res) => res.url.includes("/checkout/")
		***REMOVED***) || fail("cart submit failed");
	***REMOVED***);

	// The previous POST operation should get redirected to a dynamic URL that has a
	// path that looks like e.g. "/checkout/7". Later checkout steps get redirected
	// to paths that are prefixed with this dynamic string, which means we need to
	// save it in order to be able to (easily) verify that later POSTs are successful
	// and get redirected to e.g. "/checkout/7/shipping"
	let checkoutBaseURL = res.url;

	group("Checkout 3: billing details", function() ***REMOVED***
		let formID = res.body.match('name="form_id" value="(.*)"')[1];
		let formToken = res.body.match('name="form_token" value="(.*)"')[1];
		let formBuildID = res.body.match('name="form_build_id" value="(.*)"')[1];
		// try without setting Referer
		let formdata = ***REMOVED***
			"customer_profile_billing[commerce_customer_address][und][0][country]": "SE",
			"customer_profile_billing[commerce_customer_address][und][0][name_line]": "Mr Test",
			"customer_profile_billing[commerce_customer_address][und][0][thoroughfare]": "Gotgatan 14",
			"customer_profile_billing[commerce_customer_address][und][0][premise]": "",
			"customer_profile_billing[commerce_customer_address][und][0][postal_code]": "11846",
			"customer_profile_billing[commerce_customer_address][und][0][locality]": "Stockholm",
			"customer_profile_shipping[commerce_customer_profile_copy]": "1",
			"form_build_id": formBuildID,
			"form_token": formToken,
			"form_id": formID,
			"op": "Continue to next step"
		***REMOVED***;
		let headers = ***REMOVED*** "Content-Type": "application/x-www-form-urlencoded" ***REMOVED***;
		res = http.post(checkoutBaseURL, formdata, ***REMOVED*** headers: headers ***REMOVED***);
		check(res, ***REMOVED***
			"billing details succeeded": (res) => res.url === (checkoutBaseURL + "/shipping")
		***REMOVED***) || fail("billing details failed"); 
	***REMOVED***);

	group("Checkout 4: shipping options", function() ***REMOVED***
		let formID = res.body.match('name="form_id" value="(.*)"')[1];
		let formToken = res.body.match('name="form_token" value="(.*)"')[1];
		let formBuildID = res.body.match('name="form_build_id" value="(.*)"')[1];
		let formdata = ***REMOVED***
			"commerce_shipping[shipping_service]": "express_shipping",
			"form_build_id": formBuildID,
			"form_token": formToken,
			"form_id": formID,
			"op": "Continue to next step"
		***REMOVED***;
		let headers = ***REMOVED*** "Content-Type": "application/x-www-form-urlencoded" ***REMOVED***;
		res = http.post(checkoutBaseURL + "/shipping", formdata, ***REMOVED*** headers: headers ***REMOVED***);
		check(res, ***REMOVED***
			"shipping options succeeded": (res) => res.url === (checkoutBaseURL + "/review")
		***REMOVED***) || console.log("Select shipping failed!");
	***REMOVED***);
	
	group("Checkout 5: review and submit", function() ***REMOVED***
		let formID = res.body.match('name="form_id" value="(.*)"')[1];
		let formToken = res.body.match('name="form_token" value="(.*)"')[1];
		let formBuildID = res.body.match('name="form_build_id" value="(.*)"')[1];
		let formdata = ***REMOVED***
			"commerce_payment[payment_method]": "commerce_payment_example|commerce_payment_commerce_payment_example",
			"commerce_payment[payment_details][credit_card][number]": "4111111111111111",
			"commerce_payment[payment_details][credit_card][exp_month]": "03",
			"commerce_payment[payment_details][credit_card][exp_year]": "2019",
			"form_build_id": formBuildID,
			"form_token": formToken,
			"form_id": formID,
			"op": "Continue to next step"
		***REMOVED***;
		let headers = ***REMOVED*** "Content-Type": "application/x-www-form-urlencoded" ***REMOVED***;
		res = http.post(checkoutBaseURL + "/review", formdata, ***REMOVED*** headers: headers ***REMOVED***);
		// if this POST succeeds, it will redirect to e.g. /checkout/7/payment
		// /checkout/7/payment, in turn, will redirect to /checkout/7/paypal_ec
		// /checkout/7/paypal_ec, in turn, will redirect to /checkout/7/complete
		check(res, ***REMOVED***
			"Checkout 6: checkout complete": (res) => res.html("h1").text() === "Checkout complete"
		***REMOVED***) || fail("review and submit failed");
	***REMOVED***);
***REMOVED***

// Log out the user
function doLogout() ***REMOVED***
	check(http.get(baseURL + "/user/logout"), ***REMOVED***
		"logout succeeded": (res) => res.body.includes('<a href="/user/login">Log in')
	***REMOVED***) || fail("logout failed");
***REMOVED***

// Static resources to be loaded once per VU iteration
const staticAssets = [
	baseURL + "/modules/system/system.base.css?olqap9",
	baseURL + "/modules/system/system.menus.css?olqap9",
	baseURL + "/modules/system/system.messages.css?olqap9",
	baseURL + "/modules/system/system.theme.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/cloud_zoom/css/cloud_zoom.css?olqap9",
	baseURL + "/misc/ui/jquery.ui.core.css?olqap9",
	baseURL + "/misc/ui/jquery.ui.theme.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/libraries/jquery_ui_spinner/ui.spinner.css?olqap9",
	baseURL + "/modules/comment/comment.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/commerce_add_to_cart_confirmation/css/commerce_add_to_cart_confirmation.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/commerce_kickstart/commerce_kickstart_menus/commerce_kickstart_menus.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/date/date_api/date.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/date/date_popup/themes/datepicker.1.7.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/fences/field.css?olqap9",
	baseURL + "/modules/node/node.css?olqap9",
	baseURL + "/modules/user/user.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/views/css/views.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/ctools/css/ctools.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/commerce/modules/line_item/theme/commerce_line_item.theme.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/commerce/modules/product/theme/commerce_product.theme.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/commerce_fancy_attributes/commerce_fancy_attributes.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/alpha/css/alpha-reset.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/alpha/css/alpha-mobile.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/alpha/css/alpha-alpha.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/omega/css/formalize.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/omega/css/omega-text.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/omega/css/omega-branding.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/omega/css/omega-menu.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/omega/css/omega-forms.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/css/global.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/css/commerce_kickstart_style.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/css/omega-kickstart-alpha-default.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/css/omega-kickstart-alpha-default-narrow.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/css/commerce-kickstart-theme-alpha-default.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/css/commerce-kickstart-theme-alpha-default-narrow.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/alpha/css/grid/alpha_default/narrow/alpha-default-narrow-24.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/css/omega-kickstart-alpha-default-normal.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/css/commerce-kickstart-theme-alpha-default-normal.css?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/alpha/css/grid/alpha_default/normal/alpha-default-normal-24.css?olqap9",
	baseURL + "/misc/jquery.js?v=1.4.4",
	baseURL + "/misc/jquery.once.js?v=1.2",
	baseURL + "/misc/drupal.js?olqap9",
	baseURL + "/misc/ui/jquery.ui.core.min.js?v=1.8.7",
	baseURL + "/misc/ui/jquery.ui.widget.min.js?v=1.8.7",
	baseURL + "/profiles/commerce_kickstart/libraries/cloud-zoom/cloud-zoom.1.0.3.min.js?v=1.0.3",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/cloud_zoom/js/cloud_zoom.js?v=1.0.3",
	baseURL + "/profiles/commerce_kickstart/libraries/jquery_expander/jquery.expander.min.js?v=1.4.2",
	baseURL + "/profiles/commerce_kickstart/libraries/jquery_ui_spinner/ui.spinner.min.js?v=1.8",
	baseURL + "/profiles/commerce_kickstart/libraries/selectnav.js/selectnav.min.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/commerce_add_to_cart_confirmation/js/commerce_add_to_cart_confirmation.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/commerce_kickstart/commerce_kickstart_search/commerce_kickstart_search.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/service_links/js/twitter_button.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/service_links/js/facebook_like.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/service_links/js/google_plus_one.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/contrib/commerce_fancy_attributes/commerce_fancy_attributes.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/modules/commerce_kickstart/commerce_kickstart_product_ui/commerce_kickstart_product_ui.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/js/omega_kickstart.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/omega/js/jquery.formalize.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega/omega/js/omega-mediaqueries.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/js/commerce_kickstart_theme_custom.js?olqap9",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/logo.png",
	baseURL + "/sites/default/files/styles/product_full/public/messenger-1v1.jpg?itok=hPe-GkYY",
	baseURL + "/sites/default/files/styles/product_thumbnail/public/messenger-1v1.jpg?itok=cXkqMlMc",
	baseURL + "/sites/default/files/styles/product_thumbnail/public/messenger-1v2.jpg?itok=yyhLIuCD",
	baseURL + "/sites/default/files/styles/product_thumbnail/public/messenger-1v3.jpg?itok=uQsNvRiQ",
	baseURL + "/sites/default/files/styles/product_thumbnail/public/messenger-1v4.jpg?itok=ns9kHz1T",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/images/bg.png",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/images/picto_cart.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/picto_magnifying_glass.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/bg_product_attributes_bottom.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/bg_product_attributes_top.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/bg_add_to_cart.png",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/images/bg_block_footer_title.png",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/images/icon_facebook.png",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/images/icon_twitter.png",
	baseURL + "/profiles/commerce_kickstart/themes/commerce_kickstart_theme/images/icon_pinterest.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/picto_mastercard.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/picto_paypal.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/picto_visa_premier.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/picto_american_express.png",
	baseURL + "/misc/ui/images/ui-bg_glass_75_e6e6e6_1x400.png",
	baseURL + "/misc/ui/images/ui-icons_888888_256x240.png",
	baseURL + "/profiles/commerce_kickstart/themes/contrib/omega_kickstart/images/btn_read_more.png",
	baseURL + "/sites/default/files/messenger-1v1.jpg",
	baseURL + "/profiles/commerce_kickstart/libraries/cloud-zoom/blank.png",
];