# Micro-frontend-router
Router based on angular router which would help in loading multiple web-component based micro-frontends inside a shell application.

This is an angular library which you can install in your micro-application and instead of routerModules.forRoot()
and routerModules.forChild() you can use RoutingModule.withRoutes() provided by this library.

Your Shell-Application can still use the angular router module.

If you use angular routerModule and navigate between micro-apps in the context of the driver. You will see errors in console saying 'No matching router segment' but if you use this module those errors will be gone cause it uses a central micro-app location strategy to navigate between the pages and each time it does so, it updates the history.

Still Pending:

1. There are a few items which are still pending like the browser refresh button would set the page to the first page of the route. It is solvable but for now it's in the backlog. Will fix it soon enough

2. Navigating from microapp A ---> microApp B then cancelling the navigation would take you to microApp A. If you do imperative routing again to microApp B and if micro-appB initial route is (e.g /add-payee/find-payee). The url would only point to (/add-payee).

3. Extensive testing needed for fragments, outlets other techniques. Feel free to raise an issue if you see anything like that.

4. The routing module uses the SystemJsNgModuleLoader class which is deprecated. Need to find a better solution or write my own class.

5. Unit Test cases needs to be added.


NOTE: Doesn't support HashLocationStrategy, Implementation native to PathLocationStrategy.
