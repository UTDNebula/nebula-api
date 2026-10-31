package schema

const ClubsClubJSONB = `jsonb_build_object(
    'id', club.id,
    'slug', club.slug,
    'name', club.name,
    'alias', club.alias,
    'founding_date', (club.founding_date AT TIME ZONE 'UTC'),
    'updated_at', (club.updated_at AT TIME ZONE 'UTC'),
    'description', club.description,
    'tags', club.tags,
    'profile_image', club.profile_image,
    'banner_image', club.banner_image,
    'schools', club.schools,
    'officers', officers,
    'contacts', contacts,
    'membership_forms', membership_forms
)`

const ClubsOfficerJSONB = `jsonb_build_object(
    'name', officers.name,
    'position', officers.position
)`

const ClubsContactJSONB = `jsonb_build_object(
    'platform', contacts.platform,
    'url', contacts.url
)`

const ClubsMembershipFormJSONB = `jsonb_build_object(
    'name', membership_forms.name,
    'url', membership_forms.url
)`

const ClubsEventJSONB = `jsonb_build_object(
    'id', events.id,
    'club_id', events.club_id,
    'name', events.name,
    'description', events.description,
    'start_time', (events.start_time AT TIME ZONE 'UTC'),
    'end_time', (events.end_time AT TIME ZONE 'UTC'),
    'location', events.location,
    'image', events.image,
    'created_at', (events.created_at AT TIME ZONE 'UTC'),
    'updated_at', (events.updated_at AT TIME ZONE 'UTC')
)`
