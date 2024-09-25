# certsuite-overview
CertSuite is a dashboard that consolidates data from Quay (image pulls), DCI (test suite runs), and the CertSuite Collector. It automatically updates a centralized MySQL database, offering real-time insights into the usage of CertSuite tools and services. By providing a single, unified view of key metrics, it simplifies monitoring and reporting for stakeholders.

# Key Features
1. Quay Image Pull Tracking
   CertSuite tracks and consolidates image pull data from Quay to monitor repository usage over time.

2. DCI Test Suite Run Tracking
   Partners running CertSuite via DCI (Distributed CI) will have their test suite executions tracked.

3. CertSuite Collector Data Visualization
- Data collected from the CertSuite Collector—which includes detailed test metrics, execution logs, and other critical statistics—will be visualized.
- This data, which is already displayed on Grafana dashboards, will now also be integrated into the CertSuite dashboard for a seamless and comprehensive view, including historical tracking.

4. Automated MySQL Integration
- All data streams from Quay, DCI, and the Collector are automatically consolidated into a regularly updated MySQL database.
- This allows for easy access to all relevant metrics in one location, enabling more efficient data querying, visualization, and real-time reporting.

# Data Sources
1. Quay Image Pulls
- Tracks the number of image pulls from the CertSuite repository hosted on Quay.
- Repository: go-quay
2. DCI Test Suite Runs
- Monitors the number of CertSuite test suite executions performed by partners via DCI.
- Repository: go-dci
3. CertSuite Collector Data
- Visualizes detailed metrics on test suite executions and related performance data from the CertSuite Collector.
- Grafana Dashboard: Collector Dashboard

# Goals
The CertSuite Usage Dashboard aims to:
1. Provide a unified view of CertSuite's usage across multiple platforms.
2. Ensure that all relevant metrics are easily accessible in one place.
3. Automate data collection and reporting to reduce manual effort.

# Future Improvements
1. Integration of additional data sources for deeper insights.
2. Improved visualizations within the Google Sheet or other platforms like Grafana.
