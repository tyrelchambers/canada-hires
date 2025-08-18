# Business Demographics Reporting Feature Specification

## Overview
Community-driven system for reporting and tracking business demographics with focus on Temporary Foreign Worker (TFW) usage. Provides transparency in hiring practices to help users make informed decisions about which businesses to support.

## Core Features

### 1. User Account System
- **Account Registration**: Email link login (passwordless authentication) with email verification
- **Verification Tiers**:
  - Basic: Email verified
  - Enhanced: Phone + government ID verification
  - Trusted: Multiple verified reports + community vouching
- **Anti-Gaming**: IP tracking, email domain validation, account history

### 2. Business Directory
- **Business Profiles**: Name, address, industry, size, contact info
- **Unique Business Identification**: Prevent duplicate entries
- **Location-based Search**: City, province, postal code filtering
- **Industry Categories**: Standard business classification codes

### 3. Community Reporting System
- **Report Submission**: Authenticated users can submit TFW usage reports
- **Data Points Collected**:
  - TFW percentage estimate (0-100%)
  - Evidence type (job postings, workplace observation, LMIA records, employee testimony)
  - Hiring trends (increasing/decreasing Canadian hiring)
  - Wage practices relative to market rates
  - Working conditions assessment
  - Report confidence level (user's certainty)
  - Supporting documentation/links

### 4. Confidence Scoring Algorithm
- **Base Formula**: Each unique verified account = +1 confidence point
- **Weighted Scoring**:
  - Basic account: 1.0x weight
  - Enhanced account: 1.5x weight  
  - Trusted account: 2.0x weight
- **Visibility Threshold**: 3+ confidence points for public listing
- **Conflicting Reports**: Average TFW percentages, flag disputed entries

### 5. Business Rating System
- **Rating Categories**:
  - 🟢 **Green (Canadian-First)**: 0-20% TFW usage
  - 🟡 **Yellow (Mixed)**: 21-50% TFW usage
  - 🔴 **Red (TFW-Heavy)**: 51%+ TFW usage
  - ⚪ **Unverified**: Insufficient data or conflicting reports
- **Dynamic Updates**: Ratings update as new reports come in
- **Historical Tracking**: Track rating changes over time

### 6. Public Directory Interface
- **Search & Filter**:
  - Location (city, province, postal code)
  - Industry/business type
  - Rating category (Green/Yellow/Red)
  - Confidence level minimum
- **Business Listings**:
  - Current rating with confidence score
  - Report summary statistics
  - Alternative suggestions (Green-rated competitors in area)
  - Export functionality for boycott lists
- **Map Integration**: Visual representation of businesses by rating

### 7. Moderation & Quality Control
- **Report Validation**:
  - Automated checks for spam/duplicate content
  - Evidence requirement for high-impact claims
  - Community flagging system
- **Dispute Resolution**:
  - Business owners can contest ratings
  - Community review process for disputed reports
  - Admin override for verified false information
- **Account Monitoring**:
  - Track reporting patterns for abuse detection
  - Temporary restrictions for suspicious activity
  - Permanent bans for confirmed abuse

## Technical Architecture

### Database Schema
```
Users Table:
- id, email, verification_tier, created_at, last_active
- ip_addresses (tracking), email_domain, account_status

Businesses Table:
- id, name, address, city, province, postal_code, industry_code
- phone, website, size_category, created_at, updated_at

Reports Table:
- id, user_id, business_id, tfw_percentage, evidence_type
- confidence_level, wage_assessment, hiring_trend, conditions_rating
- supporting_docs, created_at, updated_at, is_flagged

Business_Ratings Table:
- business_id, current_rating, confidence_score, report_count
- avg_tfw_percentage, last_updated, is_disputed
```

### API Endpoints
```
Authentication:
POST /api/auth/send-login-link
GET  /api/auth/verify-login/:token
GET  /api/auth/profile

Business Management:
GET  /api/businesses (search/filter)
POST /api/businesses (create new)
GET  /api/businesses/:id
PUT  /api/businesses/:id (owner updates)

Reporting:
POST /api/reports (submit report)
GET  /api/reports/business/:id
PUT  /api/reports/:id (edit own report)
DELETE /api/reports/:id (delete own report)

Directory:
GET  /api/reports (public listings)
GET  /api/reports/export (boycott lists)
GET  /api/reports/map (location data)
```

## User Interface Components

### Report Submission Form
- Business search/selection
- TFW percentage slider (0-100%)
- Evidence type dropdown
- Supporting documentation upload
- Confidence rating (Low/Medium/High)
- Additional comments field

### Business Directory Page
- Search bar with location/industry filters
- Rating filter buttons (Green/Yellow/Red/All)
- Business cards showing rating, confidence, location
- Map view toggle
- Export to CSV/JSON functionality

### Business Detail Page
- Current rating with confidence score
- Report statistics and trends
- Recent reports (anonymized)
- Alternative businesses suggestion
- "Report this Business" button

## Implementation Priority

### Phase 1 (MVP)
1. Email link authentication (passwordless login)
2. Basic business directory
3. Simple report submission
4. Confidence scoring algorithm
5. Public business listings with ratings

### Phase 2 (Enhanced Features)  
1. Advanced search and filtering
2. Map integration
3. Enhanced user verification
4. Dispute resolution system
5. Mobile-responsive design

### Phase 3 (Advanced Features)
1. Business owner dashboard
2. Advanced analytics and trends
3. API for third-party integrations
4. Mobile app development
5. Automated LMIA data integration

## Success Metrics
- Number of verified user accounts
- Business listings with sufficient confidence scores
- Report submission rate and quality
- User engagement with directory features
- Accuracy of ratings (validated against known data)

## Legal Considerations
- User-generated content disclaimers
- Data privacy compliance (PIPEDA)
- Defamation protection measures
- Business owner right to respond
- Terms of service for report accuracy